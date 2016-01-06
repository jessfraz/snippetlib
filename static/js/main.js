var runningRequest = false;

var search = function(category, query){
    if (!runningRequest){
        runningRequest=true;
        $.ajax({
            type: "POST",
            url: "/search",
            dataType: "json",
            data: { category: category, q: query }
        }).done(function(data) { 
            $('div#results').empty();
            for (var i=0; i<data.length; i++){ 
                $('div#results').append('<a class="block" href="/'+data[i].category+'/'+data[i].slug+'">'+data[i].name+'</a>');
            }
            runningRequest=false;
        });
    }
};


$(document).ready(function(){
    prettyPrint(); 

    $('input#search-box').on('keydown', function(e){
        var $q = $('#search-box').val();

        search(category, $q);
    });

    $('input#search-box').on('keyup', function(e){
        var $q = $('#search-box').val();

        search(category, $q);
    });

    $('form.search').submit(function(e){
        e.preventDefault();
        var $q = $('#search-box').val();
        search(category, $q);
    });

    $('input[name=email]').on('keydown', function(e){
        $('input[name=email]').removeClass('error');
        $('form.mailchimp button').removeClass('error').addClass('purple');
    });

    $('form.mailchimp').submit(function(e){
        e.preventDefault();
        var email_address = $('input[name=email]').val(); 
        $('input[name=email]').removeClass('error');
        $('form.mailchimp button').removeClass('error').addClass('purple');

        // VALIDATE THE EMAIL
        $.ajax({
            url: "https://api.mailgun.net/v2/address/validate?address=" + $('form.mailchimp input[name=email]').val() + '&api_key=pubkey-9f4j0tyuxb87xk6593ba4g4ug685v340',
            dataType: 'json',
            type: "GET"
        }).done(function(response) { 
            if (response.is_valid){  // SEND THE FORM TO THE DATABASE
                $.ajax({
                    url: "/newsletter_signup",
                    type: "POST",
                    data: $('form.mailchimp').serialize()
                }).done(function(data) {
                    $('footer p.pull-right').html(':) Thanks! We will be in touch!');
                    $('form.mailchimp').remove();
                });
            } else {
                $('input[name=email]').addClass('error');
                $('form.mailchimp button').removeClass('purple').addClass('error');
            }
        });
        return false;
    });

    $('.email-scroll').on('click', function(e){
        e.preventDefault();
        $("html, body").animate({ 
            scrollTop: $(document).height()
        }, 400, function(){
            $('input[name=email]').focus();
        });
    });
});
