var is_chrome = navigator.userAgent.toLowerCase().indexOf('chrome') > -1;
if (is_chrome) {
  document.write('<link rel="stylesheet" href="/css/chrome.css" type="text/css" />');
} else {
  document.write('<link rel="stylesheet" href="/css/firefox.css" type="text/css" />');
}

$('#ssh_window_button_close').click(function() {
  var show_ssh = "false";
  $('#ssh_frame').remove();
  $('#ssh_window').hide();
});

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

$(document).ready(function() {

	parent.frames["nav_frame"].get_nav_frame_maps();

  //Hide ssh window frame
  if (show_ssh == "false") {
    
    $('#ssh_window').hide();

  } else {

    $('#ssh_window').show();

    if (selected_vm != "") {

      $('#ssh_div_container').html("<iframe name=\"ssh_frame\" id=\"ssh_frame\" style=\"height:300px;float:left;display:inline-block;frameborder:5;scrolling:yes;width:100%;resize:both;\" src=\"/ssh_frame\"></iframe>");
    }
  }

  $.ajax({
		url: '/get_vm_map/',
		type: 'post',
		dataType: 'json',
		//data : {log_path: log_path, log_position: current_log_position},
		success : function(data) {
      
      var current_id = 1
      //var button_style = "style=\"font-size : 12px; padding: 0.0em;\""        	        

			$.each(data, function(vm_name, vm_map) {
        var html_row = "<tr>"
      	
      	/*	        	
        if (vm_map["latest_ping"] == "Pass" && vm_map["latest_ssh"] == "Pass")
        {
          html_row += "<tr bgcolor=\"#00CC00\" id=\"ansible_table_row_" + vm_name + "\">\n"              
        }
        else
        {
          html_row += "<tr bgcolor=\"#FF0000\" id=\"ansible_table_row_" + vm_name + "\">\n"
        }
        */

        html_row += "<form id=\"form_" + current_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"
        
        html_row += "<td>"
        
        //html_row += "<input form=\"form_" + current_id + "\" type=\"submit\" id=\"button_docker_" + current_id + "\" class=\"custom_text_button\" name=\"button_docker\" value=\"" + vm_name + "\" >"
        html_row += "<input form=\"form_" + current_id + "\" size=\"10\" class=\"custom_link\" type=\"button\" id=\"button_vm_drilldown_" + current_id + "\" name=\"button_vm_drilldown\" value=\"" + vm_name + "\" >"

        html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vm_name"] + "\" name=\"vm_name\">\n"   
        html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + "virtualbox" + "\" name=\"platform_type\">\n"
        html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\" value=\"default_vpc\">"
        html_row += "</td>\n"
        
        html_row += "<td>"
        html_row += "<input form=\"form_" + current_id + "\" class=\"custom_text_button\" size=\"10\" type=\"text\" value=\"" + vm_map["fqdn"] + "\" name=\"fqdn_display\" readonly>"
        html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["fqdn"] + "\" name=\"fqdn\">"
        html_row += "</td>\n"
        
        html_row += "<td>"
        html_row += "<input form=\"form_" + current_id + "\" class=\"custom_text_button\" type=\"text\" value=\"" + vm_map["public_ip_address"] + "\" name=\"ip_address_display\" readonly>"
        html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_ip_address"] + "\" name=\"public_ip_address\">"
        html_row += "</td>\n"
        
        html_row += "<td>"
        html_row += "<input form=\"form_" + current_id + "\" class=\"custom_text_button\" type=\"text\" value=\"" + vm_map["OS"] + "\" name=\"OS_display\" readonly>"
        html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"
        html_row += "</td>\n"
        
        /////////////////////////////////////////////////////////////////////

        html_row += "<td>"
        html_row += "<input form=\"form_" + current_id + "\" type=\"submit\" id=\"button_launch_" + current_id + "\" class=\"power_on_button\" name=\"button_launch\" value=\"&#xF011;\" >"
        html_row += "</td>\n"

        html_row += "<td>"            
        html_row += "<input form=\"form_" + current_id + "\" type=\"submit\" id=\"button_shutdown_" + current_id + "\" class=\"power_off_button\" name=\"button_shutdown\"  value=\"&#xF011;\" >"
        html_row += "</td>\n"

        html_row += "<td>"         
        html_row += "<input form=\"form_" + current_id + "\" type=\"image\" src=\"/css/restart.png\" id=\"button_reboot_" + current_id + "\" class=\"button_image\" name=\"button_reboot\"  value=\"Reboot\" >"
        html_row += "</td>\n"	          	         	          	          
        
        //html_row += "<td>"
        //html_row += "<input form=\"form_" + current_id + "\" type=\"submit\" id=\"button_ssh_" + current_id + "\" class=\"custom_text_button\" name=\"button_ssh\" value=\"Test SSH\" >"
        //html_row += "</td>\n"	          	         
        
        html_row += "<td>"
        html_row += "<input form=\"form_" + current_id + "\" type=\"submit\" id=\"button_sshlaunch_" + current_id + "\" class=\"custom_text_button\" name=\"button_sshlaunch\" value=\"SSH\" >"
        html_row += "</td>\n"

        html_row += "<td>"            
        html_row += "<input form=\"form_" + current_id + "\" type=\"submit\" id=\"button_delete_" + current_id + "\" class=\"custom_text_button\" name=\"button_delete\"  value=\"X\" >"
        html_row += "</td>\n"
        
        html_row += "</form>\n"
        html_row += "</tr>\n"
        
        $('#table_body').append(html_row);
      
        $("#button_launch_" +  current_id).click(function(event) {
        	//event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/configure/create_vm/";
          //document.getElementById(form_id).submit();
        });         
        
        $("#button_delete_" +  current_id).click(function(event) {
        	//event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/virtualbox/";
          document.getElementById(form_id).target = ""; 
          document.getElementById(form_id).submit();
        });

        $("#button_shutdown_" +  current_id).click(function(event) {
        	//event.preventDefault();
        	
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/ansible/run_ansible/";
          //document.getElementById(form_id).submit();
        });
        
        $("#button_reboot_" +  current_id).click(function(event) {
        	
        	//event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/ansible/run_ansible/";
          //document.getElementById(form_id).submit();
        });
        
        $("#button_ssh_" +  current_id).click(function(event) {
        	//event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/virtualbox/";
          document.getElementById(form_id).target = ""; 
          document.getElementById(form_id).submit();
          //event.preventDefault();
        });
        
        $("#button_sshlaunch_" +  current_id).click(function(event) {        	
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/virtualbox/";
          document.getElementById(form_id).target = ""; 
          document.getElementById(form_id).submit();

          //event.preventDefault();
        });                   

        $("#button_vm_drilldown_" +  current_id).click(function(event) {
        	event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).target = ""; 
          document.getElementById(form_id).action = "/specific_vm/";
          document.getElementById(form_id).submit();
        });

        $("#button_docker_" +  current_id).click(function(event) {
        	event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/docker/";
          document.getElementById(form_id).submit();
        });    
        
        $("#ansible_deploy_" +  current_id).click(function(event) {
        	event.preventDefault();
          var form_id = $(this).attr("form");
          document.getElementById(form_id).action = "/ansible/run_ansible/";
          document.getElementById(form_id).submit();
        });
        
        current_id += 1
      }); 
        
		},
  }); 
  
  ///////////////////////////////////////////////////////////////////////////////////////////////////////
  /*
  $.ajax({
    url: '/get_vm_map/',
    type: 'post',
    dataType: 'json',
    //data : {log_path: log_path, log_position: current_log_position},
    success : function(data) {
    
      var current_id = 1
      var button_style = "style=\"font-size : 11px; padding: 0.4em;\""
      button_style = ""
      var html_row = ""
      var vm_name = ""
      
      html_row += "<tr id=\"vm_table_row\">"
      
      html_row += "<td>"
      html_row += "<select name=\"vm_list_selection\" id=\"vm_list_selection\">"
      
      html_row += "<option value=\"null\">&lt;Select&gt;</option>"
       
      $.each(data, function(vm_name, vm_map) {          
        html_row += "<option value=\"" + vm_map["vm_name"] + "\">" + vm_map["vm_name"] + "</option>"
      });
      
      html_row += "</select>"
      html_row += "</td>"
      
      html_row += "<td>"      
      html_row += "<div id=\"ip_address_" + current_id + "\"></div>"        
      html_row += "</td>"
              
      html_row += "<td>"
      html_row += "<input type=\"submit\" style=\"padding:1px 5px;\" name=\"Launch Nexus\" id =\"button_launch_nexus\" value=\"Launch Nexus\">"
      html_row += "</td>"
      
      html_row += "</tr>"

      $('#service_selector_table').append(html_row);
      
      $("#vm_list_selection").change(function() {
      
        var selection = $("#vm_list_selection option:selected").text();
        var ip_address = data[selection]["ip_address"]
        
        $("#ip_address_" + current_id).html(data[selection]["ip_address"]);
        $("#ip_address").val(data[selection]["ip_address"]);
      });
      
      $("#button_launch_nexus").click(function() {
      
        var selection = $("#vm_list_selection option:selected").text();
        var ip_address = data[selection]["ip_address"]
        
        window.open("https://" + ip_address + ":8443", "_blank");
      });
            
    },      
  });
  */

  ///////////////////////////////////////////////////////////////////////////////////////////////////////
  
});
