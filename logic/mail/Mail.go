package mail

import (
	"srserver/common"
	"srserver/logic/playermanager"
	"srserver/logic/utils"
	pb "srserver/proto/out"
	"time"
)

func SendSystemMail(recipient playermanager.Player, money *pb.Money) {
	mailMessage := &pb.MailMessage{
		Id:       utils.MakeId(),
		FromName: "Sys",
		FromUid:  0,
		ToUid:    0,
		Time:     time.Now().UnixMilli(),
		Title:    "ADMIN",
		Message:  "admin",
		IsReaded: false,
		IsSystem: true,
		Money:    money,
		Exp:      0,
		Fuel:     0,
		Upgrades: []*pb.CarUpgrade{},
		Items:    []*pb.Item{},
	}
	recipient.User.Mail.Mails = append(recipient.User.Mail.Mails, mailMessage)
	recipient.Save()
	if recipient.Conn != nil {
		pack := common.NewPack()
		pack.SetMethod(common.OnNewMail)
		pack.WriteProto(mailMessage)
		recipient.Conn.Write(pack.ToByteBuff())
	}

}
