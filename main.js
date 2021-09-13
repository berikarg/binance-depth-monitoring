function getDepth(){
    var symbolForm = $("#symbol").serialize();
    var limitForm = $("#limit").serialize();

    $.ajax({
      type: "POST",
      url: 'http://localhost:8080/depth',
      data: {symbol: symbolForm, limit: limitForm},
      dataType: 'json',
      success: function (depth) {
        $('#BidOrdersDepth').text(JSON.stringify(depth.bids))
        $('#BidOrdersSum').text(depth.BidsOrderSum)
        $('#AskOrdersDepth').text(JSON.stringify(depth.asks))
        $('#AskOrdersSum').text(depth.AsksOrderSum)
      },
      statusCode: {
        400: function () {
          //$('#depthContainer').style.display = 'none'
          //$('#errorContainer').style.display = 'block'
          console.log("Something")
          $('#BidOrdersDepth').text("Error: Symbol or limit are invalid")
        }
      }
    });
  }
  setInterval(getDepth, 2000)