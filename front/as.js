'use strict'

var as = new function(){
    var that = this;

    that.load = function(){
	var t_login = $('#login-templ').html();
	var t_dash = $('#dash-templ').html();
	var t_qa = $('#qa-templ').html();
	var j_as = $('#as');

	$.post('/spt/d/mg',{},function(res){
	    if(res.status != 'SUCCES'){
		j_as.html(Mustache.render(t_login));
		$('#login').on('click',function(e){
		    var mail = $('#mail').val();
		    var passwd = $('#passwd').val();

		    $.post('/spt/d/login',{
			'mail':mail,
			'passwd':passwd,
		    },function(res){
			location.reload();
		    });
		});
	    }else{
		j_as.html(Mustache.render(t_dash));
		$('a.expand').on('click',function(e){
		    var j_this = $(this);
		    
		    if(j_this.attr('toggle') != 'true'){
			j_this.attr('toggle','true');
			j_this.find('h2 > span').text('[-]');   
			j_this.siblings('div.cont').show();
		    }else{
			j_this.attr('toggle','false');
			j_this.find('h2 > span').text('[+]');
			j_this.siblings('div.cont').hide();
		    }

		    return false;
		});
	    }
	}).done(function(){
	    $.post('/spt/d/mg/qa',{},function(res){
		var i;

		console.log($('div.qa > div.cont').length);
		$('div.qa > div.cont').html(Mustache.render(t_qa,res));

		if(res.data != null){
		    for(i = 0;i < res.data.length;i++){
			$('[qaid="' + res.data[i].Id + '"]').data('qa',res.data[i]);
		    }
		}

		$('div.qa div.edit > button.submit').on('click',function(e){
		    var subject = $('div.qa div.edit > input.subject').val();
		    var clas = $('div.qa div.edit > input.clas').val();
		    var order = parseInt($('div.qa div.edit > input.order').val());
		    var body = $('div.qa div.edit > textarea').val();
		    
		    $.post('/spt/d/mg/qa_add',{
			'data':JSON.stringify({
			    'Id':$('div.qa div.edit').attr('qaid'),
			    'Subject':subject,
			    'Clas':clas,
			    'Order':order,
			    'Body':body,
			})
		    },function(res){
			location.reload();
		    });
		});
		$('div.qa div.edit > button.cancel').on('click',function(e){
		    location.reload();
		});
		$('div.qa div.list button.modify').on('click',function(e){
		    var qa = $(this).parent().data('qa');
		    $('div.qa div.edit').attr('qaid',qa.Id);
		    $('div.qa div.edit > input.subject').val(qa.Subject);
		    $('div.qa div.edit > input.clas').val(qa.Clas);
		    $('div.qa div.edit > input.order').val(qa.Order);
		    $('div.qa div.edit > textarea').val(qa.Body);
		    location.hash = "edit";
		});
		$('div.qa div.list button.delete').on('click',function(e){
		    $.post('/spt/d/mg/qa_add',{
			'data':JSON.stringify({
			    'Id':$(this).parent().attr("qaid"),
			    'Subject':"",
			})
		    },function(res){
			location.reload();
		    });
		});
	    });
	});
    };
}
