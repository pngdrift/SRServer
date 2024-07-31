package handlers

import (
	"srserver/common"
	"srserver/logic/clients"
	"srserver/logic/playermanager"
	pb "srserver/proto/out"
)

func ReadMail(pack *common.Pack, client *clients.Client) {
	id := pack.ReadLong()

	for _, mail := range client.Player.User.Mail.Mails {
		if !mail.IsReaded && mail.Id == id {
			readMail(mail, client.Player)
			break
		}
	}
}

func ReadAllMails(pack *common.Pack, client *clients.Client) {
	for _, mail := range client.Player.User.Mail.Mails {
		if !mail.IsReaded {
			readMail(mail, client.Player)
		}
	}
}

func DeleteMail(pack *common.Pack, client *clients.Client) {
	id := pack.ReadLong()
	mails := client.Player.User.Mail.Mails
	newMails := mails[:0]
	for _, mail := range mails {
		if mail.Id == id {
			if !mail.IsReaded {
				readMail(mail, client.Player)
			}
			continue
		}
		newMails = append(newMails, mail)
	}
	client.Player.User.Mail.Mails = newMails
	client.Player.Save()
}

func DeleteReadedMails(pack *common.Pack, client *clients.Client) {
	mails := client.Player.User.Mail.Mails
	unreadedMails := mails[:0]
	for _, mail := range mails {
		if !mail.IsReaded {
			unreadedMails = append(unreadedMails, mail)
		}
	}
	client.Player.User.Mail.Mails = unreadedMails
	client.Player.Save()
}

func readMail(mail *pb.MailMessage, player *playermanager.Player) {
	mail.IsReaded = true
	player.User.Fuel.Addition += mail.Fuel
	player.AddExp(mail.GetExp())
	player.Deposit(mail.GetMoney())
}
