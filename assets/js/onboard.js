$(function() {
    $('#joyreadURL').val(window.location.origin)
    
    $('#sendTestEmail').click(function(e) {
        e.preventDefault()

        $('#testEmailResult').html('<i></i>')
        
        var smtpHostname = $('#smtpHostname').val()
        var smtpPort = $('#smtpPort').val()
        var smtpUsername = $('#smtpUsername').val()
        var smtpPassword = $('#smtpPassword').val()
        var smtpTestEmail = $('#testEmailAddr').val()

        var data = {
            'smtp_hostname': smtpHostname,
            'smtp_port': smtpPort,
            'smtp_username': smtpUsername,
            'smtp_password': smtpPassword,
            'smtp_test_email': smtpTestEmail
        }

        $.ajax({
            type: 'POST',
            url: '/test-email',
            data: JSON.stringify(data),
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            success: function(response) {
                var message = response.is_email_sent ? 'Email sent successfully!' : 'Email not sent'
                $('#testEmailResult').html(message)
            },
            error: function(jqXhr, textStatus, errorThrown) {
                $('#testEmailResult').html(jqXhr.responseText)
            }
        })
    })
})