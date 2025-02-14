package content

import (
	"bufio"
	"encoding/binary"
	"hash/crc32"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"srserver/common"
	"srserver/conf"
	pb "srserver/proto/out"
)

type ContentServer struct {
	patchContainer *pb.PatchContainer
}

func NewContentServer() *ContentServer {
	server := &ContentServer{}
	server.CreatePatchContainer()
	return server
}

func (s *ContentServer) Start() {
	ln, err := net.Listen("tcp", ":8993")
	if err != nil {
		log.Fatalf("Failed to start content server: %v", err)
	}
	log.Println("Content server is listening on port 8993")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Content server err: %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *ContentServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	packetLengthBuffer := make([]byte, 4)
	for {
		_, err := reader.Read(packetLengthBuffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading packet length:", err)
			}
			return
		}
		packetLength := binary.BigEndian.Uint32(packetLengthBuffer)

		packetBuffer := make([]byte, packetLength)
		_, err = reader.Read(packetBuffer)
		if err != nil {
			log.Println("Error reading packet:", err)
			return
		}

		pack := common.NewPackInstance(packetBuffer)
		s.processPacket(pack, conn)
	}
}

func (s *ContentServer) processPacket(pack *common.Pack, conn net.Conn) {
	response := common.NewResponse(pack)
	switch pack.GetMethod() {
	case common.CheckVersion:
		clientVersion := pack.ReadInt()
		if conf.TECHNICAL_WORKS {
			response.WriteInt(1)
		} else if clientVersion < conf.MIN_CLIENT_VERSION || clientVersion > conf.MAX_CLIENT_VERSION {
			response.WriteInt(-1)
		} else {
			response.WriteInt(0)
		}
	case common.GetPatchContainer:
		response.WriteProto(s.patchContainer)
		response.WriteString(conf.LOGIC_SERVER_ADDRESS)
	case common.GetPatchFile:
		pathUrl := pack.ReadString()
		filePosition := pack.ReadInt()
		size := pack.ReadInt()
		data, err := os.ReadFile(conf.ASSETS_PATH + pathUrl)
		if err != nil {
			response.SetError(true)
			response.WriteProto(&pb.GameException{
				ErrorMessage:     err.Error(),
				ErrorDescription: "get_patch_file",
				Error:            0,
				ErrorLevel:       pb.ErrorLevel_ERROR,
			})
		} else {
			remainingBytes := len(data) - int(filePosition)
			bytesToSend := min(remainingBytes, int(size), 1024*1024*2) // TODO
			block := make([]byte, bytesToSend)
			copy(block, data[filePosition:filePosition+int32(bytesToSend)])
			response.WriteBytes(block)
		}
	default:
		return
	}
	conn.Write(response.ToByteBuff())
}

func (s *ContentServer) CreatePatchContainer() {
	var patchFiles []*pb.PatchFile
	var totalSize int64
	err := filepath.Walk(conf.ASSETS_PATH, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			crc32Value := crc32.ChecksumIEEE(fileData)
			relativePath, err := filepath.Rel(conf.ASSETS_PATH, filePath)
			if err != nil {
				return err
			}
			fileSize := info.Size()
			totalSize += fileSize
			patchFiles = append(patchFiles, &pb.PatchFile{
				Path: relativePath,
				Size: fileSize,
				Hash: int64(crc32Value),
			})
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}
	s.patchContainer = &pb.PatchContainer{
		Files: patchFiles,
	}
	log.Printf("Patch container created: total files size: %vMB", totalSize/1024/1024)
}
