
function handle_message(ctx, sk, event)
    sk:SendText(ctx, "```\n" .. json(event) .. "\n```" , event.Channel)
    return true
end