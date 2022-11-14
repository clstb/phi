ERROR_PAGE = \
    """
    <!doctype html>
    <html lang="en">
    <body>

      <div style="
                 background-color: #f54242;
                 display: flex;
                 flex-direction: column;
                 align-items: center;
                 justify-content: left;
                 font-size: calc(10px + 2vmin);
                 color: black;
                 padding-bottom: 25px;">
        <h2> Root access is disabled!</h2>

        <div> Try adding ?uname=... to request </div>
      </div>
    </body>
    </html>
    """
