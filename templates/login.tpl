
<!DOCTYPE html>
<html lang="en">

  <head>
      <title>{!{.Config.Name}!} - Login</title>
      {!{template "headcontent"}!}

      <style type="text/css">
        body > .grid {
          height: 100%;
        }
        .image {
          margin-top: -100px;
        }
        .column {
          max-width: 450px;
        }
      </style>

  </head>

  <body>

    <div class="ui middle aligned center aligned grid">
      <div class="column">
        <h2 class="ui image header">
          <i class="large sign in icon" style="min-width: 40px;"></i>
          <div class="content">
            Log-in to {!{.Config.Name}!}
          </div>
        </h2>
        <form class="ui large form" id="form">
          <div class="ui stacked segment">
            <div class="field">
              <div class="ui left icon input">
                <i class="user icon"></i>
                <input type="text" name="username" placeholder="Username">
              </div>
            </div>
            <div class="field">
              <div class="ui left icon input">
                <i class="lock icon"></i>
                <input type="password" name="password" placeholder="Password">
              </div>
            </div>
            {!{if .Config.Web.RequireOTP}!}
            <div class="field">
              <div class="ui left icon input">
                <i class="code icon"></i>
                <input type="text" name="otp" placeholder="OTP Code">
              </div>
            </div>
            {!{end}!}
            <div class="ui fluid large submit button blue">Login</div>
          </div>

          <div class="ui error message"></div>

        </form>
      </div>
    </div>

    <script>
    $(document)
    .ready(function() {
      $('.ui.form').form({})
      .api({
          url: '/login',
          method : 'POST',
          serializeForm: true,
          onError: function(data) {
            $('.error.message').html('Invalid username or password.');
            $('.error.message').css('display', 'block');
          },
          onSuccess: function(data) {
            if (data.success){
              $('.error.message').css('display', 'none');
              window.location.replace('/');
            }
          }
      });
    });
    </script>

    {!{template "tailcontent"}!}

  </body>
</html>
