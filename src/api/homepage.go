package api

import (
	format "fmt"
	"github.com/valyala/fasthttp"
	"db"
)

func HomepageHanlder(ctx *fasthttp.RequestCtx){
	parsed_query := ctx.QueryArgs().GetUintOrZero("qeq")
	channel := make (chan db.BlogPost)
	go db.GetBlogPost(parsed_query, channel)
	post := <-channel
	format.Fprintf(ctx, "Blog post: %s", post)
	//format.Printf("Type of query is %d", parsed_query)
	//format.Fprintf(ctx, "Hello from homepage.\nParams: %s", ctx.QueryArgs())
}