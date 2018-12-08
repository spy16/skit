
function handle_message(ctx, sk, event)
    sk:SendText(ctx, "Hello!", event.Channel)
    return true
end