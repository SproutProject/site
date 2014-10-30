'use strict'

var as = new function(){
    var that = this;

    that.load = function(){
	var template = $('#template').html();

	$.post('/spt/d/mg/qa',{},function(res){
	    var i;

	    if(res.data != null){
		res.show = true;
	    } 
	    $('#as').html(Mustache.render(template,res));

	    if(res.data != null){
		for(i = 0;i < res.data.length;i++){
		    $('[qaid="' + res.data[i].Id + '"]').data('qa',res.data[i]);
		}
	    }

	    $('#edit > button.submit').on('click',function(e){
		var subject = $('#edit > input.subject').val();
		var clas = $('#edit > input.clas').val();
		var order = parseInt($('#edit > input.order').val());
		var body = $('#edit > textarea').val();
		
		$.post('/spt/d/mg/qa_add',{
		    'data':JSON.stringify({
			'Id':$('#edit').attr('qaid'),
			'Subject':subject,
			'Clas':clas,
			'Order':order,
			'Body':body,
		    })
		},function(res){
		    location.reload();
		});
	    });
	    $('#edit > button.cancel').on('click',function(e){
		location.reload();
	    });
	    $('#list button.modify').on('click',function(e){
		var qa = $(this).parent().data('qa');
		$('#edit').attr('qaid',qa.Id);
		$('#edit > input.subject').val(qa.Subject);
		$('#edit > input.clas').val(qa.Clas);
		$('#edit > input.order').val(qa.Order);
		$('#edit > textarea').val(qa.Body);
		location.hash = "edit";
	    });
	    $('#list button.delete').on('click',function(e){
		$.post('/spt/d/mg/qa_add',{
		    'data':JSON.stringify({
			'Id':$(this).parent().attr("qaid"),
			'Subject':"",
		    })
		},function(res){
		    location.reload();
		});
	    });

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
	});
    };
}
