// Copyright: (c) 2019, Nathan Bliss / Sofware Automation Solutions Ltd
// GNU General Public License v2

$(document).ready(function() {		

	//Hide ssh window frame
	if (show_ssh == "false") {
	  
	  $('#ssh_window').hide();

	} else {

	  $('#ssh_window').show();

	  if (selected_vm != "") {

	    $('#ssh_div_container').html("<iframe name=\"ssh_frame\" id=\"ssh_frame\" style=\"height:300px;float:left;display:inline-block;frameborder:5;scrolling:yes;width:100%;resize:both;\" src=\"/ssh_frame\"></iframe>");
	  }
	}

	html_string = ""

	current_form_id = "master_select_form";

	html_string = "<form id=\"" + current_form_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"
	html_string += "<select id=\"" + current_form_id + "_select\" form=\"" + current_form_id + "\" name=\"master_list_selection\" required>"
	html_string += "<option  value=\"\">&lt;Select&gt;</option>"

	$.ajax({
		url: "/get_master_map_ajax/",
		type: 'post',
		dataType: 'json',
		
		success : function(data) {
			
			master_map = data;

			$.each(data, function(key,val) {
					//$.each(val, function(key2,val2) {
						//html_string += key2 + " => " + val2 + "<BR>"
					//});
					html_string += "<option value=\"" + key + "\">" + key + "</option>"
			});

			html_string += "</select><BR><BR>"			
			html_string += "<input form=\"" + current_form_id + "\" type=\"submit\" name=\"get_dashboard_token\" id=\"get_dashboard_token\"  value=\"Get Dashboard Token\">&nbsp;&nbsp;&nbsp;&nbsp;"
			html_string += "<input form=\"" + current_form_id + "\" type=\"submit\" name=\"open_kubernetes_dashboard\" id=\"open_kubernetes_dashboard\"  value=\"Open Kubernetes Dashboard\">"
			html_string += "</form>\n"

			$('#master_vm_list').html(html_string);
			
			//////////////////////////////////////////////////////////////////////////////////////////

			$('#get_dashboard_token').click(function(event) {
				
				event.preventDefault();

				var form_id = $(this).attr("form");

				if ($('#' + form_id)[0].checkValidity()) {

					vm_name = $('#' + form_id + '_select').val();

					webform_obj = {	vm_name: vm_name};

					var kubernetes_dashboard_token_frame_handle = document.createElement('iframe');
					kubernetes_dashboard_token_frame_handle.setAttribute('name', 'kubernetes_dashboard_token_frame_handle');
					kubernetes_dashboard_token_frame_handle.setAttribute('id', 'kubernetes_dashboard_token_frame_handle');
					kubernetes_dashboard_token_frame_handle.setAttribute('style', 'height:300px;float:left;display:inline-block;frameborder:5;scrolling:yes;width:100%;resize:both;');
					kubernetes_dashboard_token_frame_handle.setAttribute('src', '/kubernetes_dashboard_token_popup_frame/');

					var kubernetes_dashboard_token_frame_window = kubernetes_dashboard_token_frame_handle.src;
  	
			  	window.open(kubernetes_dashboard_token_frame_window, "kubernetes" + "-" + vm_name, 'width=800,height=430,left=0,top=100,screenX=0,screenY=100')
			  	kubernetes_dashboard_token_frame_window.webform_obj = webform_obj
				}

			});

			//////////////////////////////////////////////////////////////////////////////////////////

			$('#open_kubernetes_dashboard').click(function(event) {
				
				event.preventDefault();

				var form_id = $(this).attr("form");

				if ($('#' + form_id)[0].checkValidity()) {
					
					vm_name = $('#' + form_id + '_select').val();

          document.getElementById(form_id).action = "https://" + master_map[vm_name]["public_ip_address"] + ":32321/";
          document.getElementById(form_id).target = "_blank";
          document.getElementById(form_id).submit();
				}
			});

		},
	});

});
