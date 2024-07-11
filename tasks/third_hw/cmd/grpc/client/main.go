package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/l-brawler-l/go_test/tasks/third_hw/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Command struct {
	Port   int
	Host   string
	Cmd    string
	Name   string
	NewName string
	Amount int
}

func (cmd *Command) Do(c proto.BankAccountsClient, ctx context.Context) error {
	switch cmd.Cmd {
	case "create":
		if err := cmd.create(c, ctx); err != nil {
			return fmt.Errorf("create account failed: %w", err)
		}
		return nil
	case "get":
		if err := cmd.get(c, ctx); err != nil {
			return fmt.Errorf("get account failed: %w", err)
		}

		return nil
	case "delete":
		if err := cmd.delete(c, ctx); err != nil {
			return fmt.Errorf("create account failed: %w", err)
		}
		return nil
	case "patch":
		if err := cmd.patch(c, ctx); err != nil {
			return fmt.Errorf("create account failed: %w", err)
		}
		return nil
	case "change":
		if err := cmd.change(c, ctx); err != nil {
			return fmt.Errorf("create account failed: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown command %s", cmd.Cmd)
	}
}

func (cmd *Command) create(c proto.BankAccountsClient, ctx context.Context) error {
	request := proto.CreateAccountRequest{
		Name:   cmd.Name,
		Amount: int32(cmd.Amount),
	}

	_, err := c.Create(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Command) get(c proto.BankAccountsClient, ctx context.Context) error {
	request := proto.GetAccountRequest{
		Name:   cmd.Name,
	}

	res, err := c.Get(ctx, &request)
	if err != nil {
		return err
	}
	fmt.Printf("name: %s; amount: %d", res.Name, res.Amount)
	return nil
}

func (cmd *Command) patch(c proto.BankAccountsClient, ctx context.Context) error {
	request := proto.PatchAccountRequest{
		Name:   cmd.Name,
		Amount: int32(cmd.Amount),
	}

	_, err := c.Patch(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Command) delete(c proto.BankAccountsClient, ctx context.Context) error {
	request := proto.DeleteAccountRequest{
		Name:   cmd.Name,
	}

	_, err := c.Delete(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Command) change(c proto.BankAccountsClient, ctx context.Context) error {
	request := proto.ChangeAccountRequest{
		Name:   cmd.Name,
		NewName: cmd.NewName,
	}

	_, err := c.Change(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}


func main() {
	portVal := flag.Int("port", 4567, "server port")
	hostVal := flag.String("host", "0.0.0.0", "server host")
	cmdVal := flag.String("cmd", "", "command to execute")
	nameVal := flag.String("name", "", "name of account")
	newnameVal := flag.String("newname", "", "new name of account")
	amountVal := flag.Int("amount", 0, "amount of account")

	flag.Parse()

	cmd := Command{
		Port:   *portVal,
		Host:   *hostVal,
		Cmd:    *cmdVal,
		Name:   *nameVal,
		NewName: *newnameVal,
		Amount: *amountVal,
	}


	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cmd.Host, cmd.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = conn.Close()
	}()
	
	c := proto.NewBankAccountsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := cmd.Do(c, ctx); err != nil {
		panic(err)
	}
}