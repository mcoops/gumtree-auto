package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func waitVisibleWithTimeout(sel string, timeout time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if err := chromedp.WaitVisible(sel).Do(ctx); err != nil {
			return err
		}
		return nil
	})
}

func getNodesWithTimeout(ctx context.Context, sel string, nodes *[]*cdp.Node, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Nodes(sel, nodes, chromedp.ByQueryAll)); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	log.Println("Logging into gumtree")
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.gumtree.com.au/t-login-form.html?sl=true"),
		chromedp.WaitVisible(`#login-form`, chromedp.ByID),
	); err != nil {
		log.Fatal(err)
	}

	// Fill in login form
	if err := chromedp.Run(ctx,
		chromedp.SendKeys(`#login-email`, os.Getenv("USERNAME")),
		chromedp.SendKeys(`#login-password`, os.Getenv("PASSWORD")),
	); err != nil {
		log.Fatal(err)
	}

	if err := chromedp.Run(ctx,
		chromedp.Click(`#btn-submit-login`),
		chromedp.WaitVisible(`#nav-my`, chromedp.ByID),
	); err != nil {
		log.Fatal(err)
	}

	log.Println("Logged in successfully!")

	// goto my ads
	log.Println("Gathering ads which need reposting")
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.gumtree.com.au/m-my-ads.html?c=1&size=50"),
		waitVisibleWithTimeout(`#my-adlisting`, 5*time.Second),
	); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("No ads need reposting")
			return
		}
		log.Fatal(err)
	}

	var nodes []*cdp.Node

	if err := getNodesWithTimeout(ctx, ".repost-ad-free", &nodes, 5*time.Second); err != nil {
		if err == context.DeadlineExceeded {
			log.Println("No ads required reposting")
			return
		}
		log.Fatal(err)
	}

	log.Printf("Found %d elements with class 'repost-ad-free'\n", len(nodes))

	for _, node := range nodes {
		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://www.gumtree.com.au"+node.Attributes[3]),
			chromedp.WaitVisible(`#feature-packages-wrapper`, chromedp.ByID), // confirm repost successful
		); err != nil {
			log.Fatal(err)
			continue
		}
		log.Printf("Reposted ad with URL %s\n" + node.Attributes[3])
	}
}
