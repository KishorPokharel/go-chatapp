const mh = document.querySelector('.message-history');
mh.scrollTop = mh.scrollHeight;

$(function() {
	const messageInput = $('.message-input');
	//$('.message-history > :last-child')[0].scrollIntoView(false);
	const messageHistory = $('.message-history');
	const onlinePeople = $('.online-people');
	//messageHistory.scrollTop(messageHistory[0].scrollHeight);

	let socket = null;
	window.socket = socket;
	if(!window["WebSocket"]) {
		alert("browser not supported")
	} else {
		socket = new WebSocket("ws://localhost:3000/room")
		socket.onclose = function () {
			alert("Connection closed by server")
		}
		socket.onmessage = function (e) {
			let msg = JSON.parse(e.data)
            //console.log(msg);
			if(msg.EventType === "userJoin") {
				onlinePeople.append(
					$("<div>").attr("data-username", msg.Name).addClass("person").append(
						$("<span>").addClass("online-icon"),
						$("<span>").addClass("person-username").text(msg.Name)
					)
				)
			} else if (msg.EventType === "userLeave") {
				$(`.person[data-username="${msg.Name}"`)[0].remove();
			} else if (msg.EventType === "messageSent") {
				messageHistory.append(
					$("<div>").addClass("message-thread message-thread--sent").append(
						$("<div>").addClass("message-thread__meta").append(
							$("<span>").addClass("message-thread__createdat")
                                       .text(moment(msg.CreatedAt).format('MM/DD/YYYY h:mm A'))
						),
						$("<div>").addClass("message-thread__body").text(msg.Message)
					)
				)
				$(".message-history > :last-child")[0].scrollIntoView(false);
			} else if (msg.EventType === "messageReceived") {
				messageHistory.append(
					$("<div>").addClass("message-thread message-thread--received").append(
						$("<div>").addClass("message-thread__meta").append(
							$("<span>").addClass("message-thread__username").text(msg.Name),
							$("<span>").addClass("message-thread__createdat")
                                       .text(moment(msg.CreatedAt).format('MM/DD/YYYY, h:mm A'))
						),
						$("<div>").addClass("message-thread__body").text(msg.Message)
					)
				)
				$(".message-history > :last-child")[0].scrollIntoView(false);

            }
		}
	}

	messageInput.on("keyup", (e) => {
		if (e.keyCode === 13) {
			if(!messageInput.val()) {
				return;
			}
			socket.send(JSON.stringify({Message: messageInput.val()}))
			messageInput.val("")
		}
	})
})































