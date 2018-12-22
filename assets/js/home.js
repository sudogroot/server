$(function() {
    $('#triggerAddNewBooks').on('click', function(e) {
        e.preventDefault()
        $('#addNewBooks').trigger('click')
    })

    $('#addNewBooks').on('change', function(e) {
        var data = new FormData()
        console.log(e.target.files)

        $.each(e.target.files, function(i, file) {
            console.log(file)
            data.append('upload[]', file)
        });

        $.ajax({
            type: 'POST',
            url: '/upload-books',
            data: data,
            contentType: false,
            processData: false,
            cache: false,
            contentType: false,
            success: function(response) {
                alert(response["status"])
            },
            error: function(jqXhr, textStatus, errorThrown) {
                $('#testEmailResult').html(jqXhr.responseText)
            }
        })
    })
})