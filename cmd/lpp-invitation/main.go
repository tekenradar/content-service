package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/coneno/logger"
	"github.com/influenzanet/messaging-service/pkg/api/email_client_service"
	"github.com/tekenradar/content-service/pkg/dbs/contentdb"
	"github.com/tekenradar/content-service/pkg/types"
	"google.golang.org/grpc"
)

const (
	DefaultGRPCMaxMsgSize = 4194304

	instanceID = "tekenradar"
)

var (
	dbService *contentdb.ContentDBService
)

func main() {
	conf := readConfig()

	dbService = contentdb.NewContentDBService(conf.ContentDBConfig, conf.InstanceIDs)

	if conf.RunParticipantCreation {
		slog.Info("Start LPplus participant creation process", "force", conf.ForceReplace, "csv", conf.CSVPath)
		participants := readCSVFile(conf.CSVPath, []rune(conf.Separator)[0])
		createParticipants(participants, conf.ForceReplace)
	}

	// send invitation email to participants who did not receive an invitation yet
	if conf.RunInvitationSending {
		emailClient, emailServiceClose := connectToEmailService(conf.EmailClientURL, DefaultGRPCMaxMsgSize)
		defer emailServiceClose()

		// Load invitation template
		emailTemplate, err := os.ReadFile(conf.InvitationEmailTemplatePath)
		if err != nil {
			slog.Error("Unable to read invitation email template", "error", err, slog.String("path", conf.InvitationEmailTemplatePath))
			return

		}

		uninvitedParticipnts, err := dbService.FindUninvitedLPPParticipants(instanceID)
		if err != nil {
			slog.Error("Unable to find uninvited participants", "error", err)
			return
		}

		if len(uninvitedParticipnts) == 0 {

			fmt.Println("\n=====================================================================")
			slog.Info("No uninvited participants found")
			fmt.Println("=====================================================================")
			fmt.Println()
			fmt.Println()
			return
		}

		target := len(uninvitedParticipnts)
		sent := 0

		for _, p := range uninvitedParticipnts {
			slog.Debug("Send invitation email", "pid", p.PID)

			content, err := ResolveTemplate(
				"lpp-invitation",
				string(emailTemplate),
				map[string]string{
					"pid":  p.PID,
					"name": p.ContactInfos.Name,
				},
			)
			if err != nil {
				logger.Error.Printf("invitation email could not be generated: %v", err)
				continue
			}

			_, err = emailClient.SendEmail(context.Background(), &email_client_service.SendEmailReq{
				To:      []string{p.ContactInfos.Email},
				Subject: conf.InvitationEmailSubject,
				Content: content,
			})
			if err != nil {
				slog.Error("Unable to send invitation email", "error", err, "pid", p.PID)
				continue
			}

			err = dbService.UpdateLPPParticipantInvitationSentAt(instanceID, p.PID, time.Now())
			if err != nil {
				slog.Error("Unable to update invitation sent at", "error", err, "pid", p.PID)
				continue
			}
			sent++
		}

		fmt.Println("\n=====================================================================")
		slog.Info("Sent invitation emails", "sent", sent, "target", target)
		fmt.Println("=====================================================================")
		fmt.Println()
		fmt.Println()
	}
}

func readCSVFile(filePath string, separator rune) []types.LPPParticipant {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = separator
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse CSV file: %v", err)
	}

	if len(records) < 1 {
		slog.Error("CSV file is empty")
		return []types.LPPParticipant{}
	}

	participants := []types.LPPParticipant{}

	colNames := records[0]

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 5 {
			slog.Error("CSV file entry is invalid", "record", record)
			continue
		}
		// Process each record (e.g., print it)
		p := types.LPPParticipant{
			PID: record[1],
			ContactInfos: &types.LPPParticipantContactInfos{
				Email: record[2],
				Name:  record[3],
			},
			StudyData: make(map[string]string),
		}
		for i := 4; i < len(record); i++ {
			p.StudyData[colNames[i]] = record[i]
		}
		participants = append(participants, p)
	}
	return participants
}

func createParticipants(participants []types.LPPParticipant, replace bool) {
	counter := 0
	for _, p := range participants {
		var err error
		if replace {
			err = dbService.ReplaceLPPParticipant(instanceID, p)
		} else {
			_, err = dbService.AddLPPParticipant(instanceID, p)
		}
		if err != nil {
			slog.Error("Unable to create participant", "error", err)
			continue
		}
		counter++
	}
	fmt.Println("\n=====================================================================")
	slog.Info("Created participants", "created", counter, "target", len(participants))
	fmt.Println("=====================================================================")
	fmt.Println()
	fmt.Println()
}

func connectToGRPCServer(addr string, maxMsgSize int) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxMsgSize),
		grpc.MaxCallSendMsgSize(maxMsgSize),
	))
	if err != nil {
		logger.Error.Fatalf("failed to connect to %s: %v", addr, err)
	}
	return conn
}

func connectToEmailService(addr string, maxMsgSize int) (client email_client_service.EmailClientServiceApiClient, close func() error) {
	serverConn := connectToGRPCServer(addr, maxMsgSize)
	return email_client_service.NewEmailClientServiceApiClient(serverConn), serverConn.Close
}

func ResolveTemplate(tempName string, templateDef string, contentInfos map[string]string) (content string, err error) {
	if strings.TrimSpace(templateDef) == "" {
		logger.Error.Printf("error: empty template %s", tempName)
		return "", errors.New("empty template `" + tempName)
	}
	tmpl, err := template.New(tempName).Parse(templateDef)
	if err != nil {
		logger.Error.Printf("error when parsing template %s: %v", tempName, err)
		return "", err
	}
	var tpl bytes.Buffer

	err = tmpl.Execute(&tpl, contentInfos)
	if err != nil {
		logger.Error.Printf("error when executing template %s: %v", tempName, err)
		return "", err
	}
	return tpl.String(), nil
}
