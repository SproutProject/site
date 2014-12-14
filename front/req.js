var reqform = new function(){
    var that = this;
    that.load = function(clas){
	var j_main;
	var recaptcha_id;
	var t_prepro = $('#prepro-templ').html();
	var prolist = null;

	if(clas == 0){
	    j_main = $('#reqalg');
	}else{
	    j_main = $('#reqlang');
	}

	recaptcha_id = grecaptcha.render(
		j_main.find('div.recaptcha > div.form')[0],{
		    'sitekey':'6LfuNf8SAAAAAHCcjfmarCcMKJfzXebos3TwmLRu'
		});
	
	j_main.find('div.recaptcha > button.send').on('click',function(e){
	    var j_btn = $(this);
	    verify = j_main.find('div.verify > input.verify').val();

	    j_btn.hide();
	    j_btn.siblings('span.msg').show();

	    $.post('/spt/d/req/getpre',{
		'data':JSON.stringify({
		    'Clas':clas,
		    'Recaptcha':grecaptcha.getResponse(recaptcha_id)
		})
	    },function(res){
		var i,j;

		if(res.status != 'SUCCES'){
		    alert('驗證碼錯誤');
		    location.reload();
		    return 0;
		}
		
		prolist = res.data;
		for(i = 0;i < prolist.length;i++){
		    pro = prolist[i];
		    pro.index = i + 1;
		    for(j = 0;j < pro.Option.length;j++){
			pro.Option[j] = {
			    'desc':pro.Option[j],
			    'value':j
			};
		    }
		}

		j_main.find('div.recaptcha').hide();
		j_main.find('div.prepro').show();
		j_main.find('div.prepro').html(Mustache.render(t_prepro,prolist));

		j_main.find('div.prepro > button').on('click',function(e){
		    ans = [];
		    for(i = 0;i < prolist.length;i++){
			val = j_main.find('div.prepro input[name=' + (i + 1) + ']:checked').val();
			ans[i] = parseInt(val);
		    }

		    $.post('/spt/d/req/checkpre',{
			'data':JSON.stringify(ans)
		    },function(res){
			if(res.status == 'SUCCES'){
			    j_main.find('div.prepro').hide();
			    j_main.find('div.checkmail').show();
			}else{
			    alert('答案錯誤，請重新填寫');
			    location.reload();
			}
		    });
		});
	    });
	});

	j_main.find('div.checkmail > button.send').on('click',function(e){
	    var j_btn = $(this);
	    var mail = j_main.find('div.checkmail > input.mail').val();
	    var repeat = j_main.find('div.checkmail > input.mail-repeat').val();

	    if(mail != repeat){
		alert('信箱輸入有錯誤');
	    }else{
		$.post('/spt/d/req/checkmail',{
		    'data':mail
		},function(res){
		    if(res.status == 'SUCCES'){
			j_main.find('div.checkmail').hide();
			j_main.find('div.verify').show();
		    }else{
			alert('信箱尚未符合申請資格(分數未達門檻)');
			j_btn.show();
			j_btn.siblings('span.msg').hide();
		    }
		});

		j_btn.hide();
		j_btn.siblings('span.msg').show();
	    }
	});
	j_main.find('div.verify > button.send').on('click',function(e){
	    var j_btn = $(this);
	    verify = j_main.find('div.verify > input.verify').val();

	    $.post('/spt/d/req/verify',{
		'data':verify
	    },function(res){
		if(res.status == 'SUCCES'){
		    j_main.find('div.verify').hide();
		    j_main.find('div.data').show();
		}else{
		    alert('驗證碼錯誤');
		    j_btn.show();
		    j_btn.siblings('span.msg').hide();
		}
	    });

	    j_btn.hide();
	    j_btn.siblings('span.msg').show();
	});
	j_main.find('div.data > button.send').on('click',function(e){
	    var i;
	    var j_basic = j_main.find('div.data > div.basic > input,div.data > div.basic > select,div.data > div.basic > textarea');
	    var j_textarea = j_main.find('div.data > textarea');
	    var j_btn = $(this);
	    var data = [];
	    var j_check;
	    var from = '';

	    for(i = 0;i < j_basic.length;i++){
		data.push($(j_basic[i]).val());
	    }

	    j_check = j_main.find('div.data > div.ad > input[name=from]:checked');
	    for(i = 0;i < j_check.length;i++){
		from += $(j_check[i]).val() + ',';
	    }
	    from += j_main.find('div.data > div.ad > input[type=textbox][name=from]').val();
	    data.push(from);

	    for(i = 0;i < j_textarea.length;i++){
		data.push($(j_textarea[i]).val());
	    }

	    $.post('/spt/d/req/data',{
		'data':JSON.stringify(data)
	    },function(res){
		if(res.status == 'SUCCES'){
		    j_main.find('div.data').hide();
		    j_main.find('div.done').show();
		}
	    });
	});
    };
};
