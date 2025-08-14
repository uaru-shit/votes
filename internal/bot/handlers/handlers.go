package handlers

import (
	"fmt"
	"time"

	"github.com/uaru-shit/votes/internal/domain"
	"github.com/uaru-shit/votes/pkg/utils"
	tb "gopkg.in/telebot.v4"
)

func handleBan(ctx domain.Context, msg *tb.Message, member *tb.ChatMember) {
	bot := ctx.Bot()
	log := ctx.Log()

	if err := bot.Ban(ctx.Chat(), member); err != nil {
		log.Error("cannot ban user", utils.ErrorAttr(err))

		if _, err := bot.Reply(msg, "Чота не могу забанить"); err != nil {
			log.Error("can't even cry", utils.ErrorAttr(err))
		}
		return
	}

	if _, err := bot.Reply(msg, "Ban nyuuu"); err != nil {
		log.Error("failed to reply to poll", utils.ErrorAttr(err))
	}
}

func handleUnban(ctx domain.Context, msg *tb.Message, member *tb.ChatMember) {
	bot := ctx.Bot()
	log := ctx.Log()

	if err := bot.Unban(ctx.Chat(), member.User); err != nil {
		log.Error("cannot unban user", utils.ErrorAttr(err))

		if _, err := bot.Reply(msg, "Чота не могу разбанить"); err != nil {
			log.Error("can't even cry", utils.ErrorAttr(err))
		}
		return
	}

	if _, err := bot.Reply(msg, "Разбан nyuuu"); err != nil {
		log.Error("failed to reply to poll", utils.ErrorAttr(err))
	}
}

func handlePollResults(ctx domain.Context, msg *tb.Message, member *tb.ChatMember) {
	time.Sleep(time.Minute)

	bot := ctx.Bot()
	log := ctx.Log()

	poll, err := bot.StopPoll(msg)
	if err != nil {
		log.Error("failed to stop the poll", utils.ErrorAttr(err))

		return
	}

	if poll.Options[0].VoterCount > poll.Options[1].VoterCount {
		handleBan(ctx, msg, member)
	} else {
		handleUnban(ctx, msg, member)
	}
}

func HandleVoteban(ctx domain.Context) error {
	bot := ctx.Bot()

	if !ctx.Message().FromGroup() {
		return ctx.Reply("В лс не баню сори")
	}

	if ctx.Message().ReplyTo == nil {
		return ctx.Reply("Ответь на сообщение кого забанить")
	}

	userToBan := ctx.Message().ReplyTo.Sender
	member, err := bot.ChatMemberOf(ctx.Chat(), userToBan)

	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	admins, err := bot.AdminsOf(ctx.Chat())
	if err != nil {
		return fmt.Errorf("failed to get admins: %w", err)
	}

	if utils.IsAdmin(userToBan.ID, admins) {
		return ctx.Reply("ммм не")
	}

	if !utils.BotCanMute(ctx.BotUser().ID, admins) {
		return ctx.Reply("Админом меня сделай, олух")
	}

	msg, err := bot.Reply(ctx.Message(), &tb.Poll{
		Question:  "Забанить или разбанить?",
		Anonymous: false,
		Options: []tb.PollOption{
			{Text: "Бан"},
			{Text: "Разбан"},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send poll: %w", err)
	}

	go handlePollResults(ctx, msg, member)

	return nil
}
