
<!DOCTYPE html>
<html lang="en">

  <head>
      <title>OTP Generation Utility</title>
      {!{template "headcontent"}!}

      <style type="text/css">

      </style>

  </head>

  <body>

    <div class="ui middle aligned center aligned grid">
      <div class="column">
        <h2 class="ui image header">
          <i class="large code icon" style="min-width: 40px;"></i>
          <div class="content">
            OTP Generation Utility
          </div>
        </h2>
        <div>
          <table class="ui celled table">
            <thead>
              <tr><th>Field</th>
              <th>Value</th>
            </tr></thead>
            <tbody>
              <tr>
                <td>Enrollment QR</td>
                <td><img class="ui medium image" src="data:image/png;base64,{!{.QR_DATA}!}" /></td>
              </tr>
              <tr>
                <td>Issuer</td>
                <td>{!{.Key.Issuer}!}</td>
              </tr>
              <tr>
                <td>Account Name</td>
                <td>{!{.Key.AccountName}!}</td>
              </tr>
              <tr>
                <td>Secret</td>
                <td>{!{.Key.Secret}!}</td>
              </tr>
              <tr>
                <td>URL Encoded</td>
                <td>{!{.Key.String}!}</td>
              </tr>
            </tbody>
          </table>

        </div>
      </div>
    </div>

    {!{template "tailcontent"}!}

  </body>
</html>
