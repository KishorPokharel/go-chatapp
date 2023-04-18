const isProd = false;
const mh = document.querySelector('.message-history');
mh.scrollTop = mh.scrollHeight;

function trace(v) {
  if (!isProd) {
    console.log(v);
  }
}

$(function () {
  const username = window.localStorage.getItem('user');
  const messageInput = $('.message-input');
  //$('.message-history > :last-child')[0].scrollIntoView(false);
  const messageHistory = $('.message-history');
  const onlinePeople = $('.online-people');
  //messageHistory.scrollTop(messageHistory[0].scrollHeight);

  let socket = null;
  window.socket = socket;
  if (!window['WebSocket']) {
    alert('browser not supported');
  } else {
    let protocol = window.location.protocol === 'http:' ? 'ws:' : 'wss:';
    socket = new WebSocket(`${protocol}//localhost:3000/room`);
    socket.onclose = function () {
      alert('Connection closed by server');
    };
    socket.onmessage = function (e) {
      let msg = JSON.parse(e.data);
      trace(msg);
      if (msg.type === 'join') {
        onlinePeople.empty();
        for (let name of msg.clients) {
          onlinePeople.append(
            $('<div>')
              .attr('data-username', name)
              .addClass('person')
              .append(
                $('<span>').addClass('online-icon'),
                $('<span>').addClass('person-username').text(name)
              )
          );
        }
      } else if (msg.type === 'leave') {
        $(`.person[data-username="${msg.sender}"`)[0].remove();
      } else if (msg.type === 'message') {
        if (username != msg.sender) {
          messageHistory.append(
            $('<div>')
              .addClass('message-thread message-thread--received')
              .append(
                $('<div>')
                  .addClass('message-thread__meta')
                  .append(
                    $('<span>')
                      .addClass('message-thread__username')
                      .text(msg.sender),
                    $('<span>')
                      .addClass('message-thread__createdat')
                      .text(
                        moment(msg.date).format('MM/DD/YYYY HH:mm A', {
                          trim: false,
                          useGrouping: false,
                        })
                      )
                  ),
                $('<div>').addClass('message-thread__body').text(msg.content)
              )
          );
          $('.message-history > :last-child')[0].scrollIntoView(false); // scroll to bottom
        } else if (username === msg.sender) {
          messageHistory.append(
            $('<div>')
              .addClass('message-thread message-thread--sent')
              .append(
                $('<div>')
                  .addClass('message-thread__meta')
                  .append(
                    $('<span>')
                      .addClass('message-thread__createdat')
                      .text(
                        moment(msg.date).format('MM/DD/YYYY HH:mm A', {
                          trim: false,
                          useGrouping: false,
                        })
                      )
                  ),
                $('<div>').addClass('message-thread__body').text(msg.content)
              )
          );
          $('.message-history > :last-child')[0].scrollIntoView(false); // scroll to bottom
        }
      }
    };
  }

  messageInput.on('keyup', (e) => {
    if (e.keyCode === 13) {
      if (!messageInput.val()) {
        return;
      }
      socket.send(
        JSON.stringify({ type: 'message', content: messageInput.val() })
      );
      messageInput.val('');
    }
  });
});
