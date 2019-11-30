function get_nav_frame_maps() {

	$('#virtualbox_vm_list').html("");
	$('#aws_managed_vm_list').html("");
	$('#gcp_managed_vm_list').html("");
	$('#aws_imported_vm_list').html("");
	$('#gcp_imported_vm_list').html("");

	var current_id = 1

	$.ajax({
		url: '/get_vm_map/',
		type: 'post',
		dataType: 'json',
		
		success : function(data) {      	        
	        
			$.each(data, function(vm_name, vm_map) {

				var html_row = ""
				         
				html_row += "<form id=\"form_" + current_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n";	          	          
				
				html_row += "<li form=\"form_" + current_id + "\" id=\"button_vm_drilldown_" + current_id + "\" name=\"button_vm_drilldown\"><a href=\"/virtualbox/\">" + vm_map["vm_name"] + "</a></li>"

				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vm_name"] + "\" name=\"vm_name\">\n"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + "virtualbox" + "\" name=\"platform_type\">\n"   
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["fqdn"] + "\" name=\"fqdn\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["ip_address"] + "\" name=\"ip_address\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"			

				html_row += "</form>\n"	        	     

				$('#virtualbox_vm_list').append(html_row);

				///////////////////////////////////////////////////////////////////////////

				$("#button_vm_drilldown_" +  current_id).click(function(event) {
					event.preventDefault();					
					var form_id = $(this).attr("form");
					document.getElementById(form_id).target = "operations_frame"; 
					document.getElementById(form_id).action = "/specific_vm/";
					document.getElementById(form_id).submit();
				});

     		current_id += 1     
        }); 	      
       },
	});
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	$.ajax({
		url: '/get_aws_map/',
		type: 'post',
		dataType: 'json',
		
		success : function(data) {	        			

			$.each(data, function(vm_name, vm_map) {

				var html_row = ""	          	          
				
				html_row += "<form id=\"form_" + current_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"

				html_row += "<li form=\"form_" + current_id + "\" id=\"button_vm_drilldown_" + current_id + "\" name=\"button_vm_drilldown\"><a href=\"/aws/\">" + vm_map["vm_name"] + "</a></li>"

				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vm_name"] + "\" name=\"vm_name\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + "aws" + "\" name=\"platform_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vpc_id"] + "\" name=\"vpc_id\" value=\"default_vpc\">"					          	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"	          	         
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["region"] + "\" name=\"region\">"	          	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_dns"] + "\" name=\"public_dns\">"	          	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_ip_address"] + "\" name=\"public_ip_address\">"	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["private_ip_address"] + "\" name=\"private_ip_address\">"

				html_row += "</form>\n"	          

				$('#aws_managed_vm_list').append(html_row);            

				$("#button_vm_drilldown_" +  current_id).click(function(event) {
					event.preventDefault();					
					var form_id = $(this).attr("form");
					document.getElementById(form_id).target = "operations_frame"; 
					document.getElementById(form_id).action = "/specific_vm/";
					document.getElementById(form_id).submit();
				});

			current_id += 1	          	          
	    });	       
		},
	});
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	$.ajax({
		url: '/get_aws_imported_map/',
		type: 'post',
		dataType: 'json',
		
		success : function(data) {	        			

			$.each(data, function(vm_name, vm_map) {

				var html_row = ""	          	          
				
				html_row += "<form id=\"form_" + current_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"

				html_row += "<li form=\"form_" + current_id + "\" id=\"button_vm_drilldown_" + current_id + "\" name=\"button_vm_drilldown\"><a href=\"/aws/\">" + vm_map["vm_name"] + "</a></li>"

				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vm_name"] + "\" name=\"vm_name\">"				
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + "aws" + "\" name=\"platform_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vpc_id"] + "\" name=\"vpc_id\" value=\"default_vpc\">"					          	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"	          	         
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["region"] + "\" name=\"region\">"	          	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_dns"] + "\" name=\"public_dns\">"	          	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_ip_address"] + "\" name=\"public_ip_address\">"	          
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["private_ip_address"] + "\" name=\"private_ip_address\">"	           

				html_row += "</form>\n"	          

				$('#aws_imported_vm_list').append(html_row);            

				$("#button_vm_drilldown_" +  current_id).click(function(event) {
					event.preventDefault();					
					var form_id = $(this).attr("form");
					document.getElementById(form_id).target = "operations_frame"; 
					document.getElementById(form_id).action = "/specific_vm/";
					document.getElementById(form_id).submit();
				});

			current_id += 1	          	          
	    });	       
		},
	});
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	$.ajax({
		url: '/get_gcp_map/',
		type: 'post',
		dataType: 'json',
		
		success : function(data) {          			

			$.each(data, function(vm_name, vm_map) {

				var html_row = ""
				           
				html_row += "<form id=\"form_" + current_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"                        

				html_row += "<li form=\"form_" + current_id + "\" id=\"button_vm_drilldown_" + current_id + "\" name=\"button_vm_drilldown\"><a href=\"/gcp/\">" + vm_map["vm_name"] + "</a></li>"

				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vm_name"] + "\" name=\"vm_name\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + "gcp" + "\" name=\"platform_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vpc_id"] + "\" name=\"vpc_id\" value=\"default_vpc\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"                        
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["region"] + "\" name=\"region\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["project_id"] + "\" name=\"project_id\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["fqdn"] + "\" name=\"fqdn\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["userid"] + "\" name=\"userid\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_ip_address"] + "\" name=\"public_ip_address\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["private_ip_address"] + "\" name=\"private_ip_address\">"                         

				html_row += "</form>\n"

				$('#gcp_managed_vm_list').append(html_row);            

				$("#button_vm_drilldown_" +  current_id).click(function(event) {
					event.preventDefault();					
					var form_id = $(this).attr("form");
					document.getElementById(form_id).target = "operations_frame"; 
					document.getElementById(form_id).action = "/specific_vm/";
					document.getElementById(form_id).submit();
				});

				current_id += 1
			}); 
		},
  });
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	$.ajax({
		url: '/get_gcp_imported_map/',
		type: 'post',
		dataType: 'json',
		
		success : function(data) {          			

			$.each(data, function(vm_name, vm_map) {

				var html_row = ""
				           
				html_row += "<form id=\"form_" + current_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"                        

				html_row += "<li form=\"form_" + current_id + "\" id=\"button_vm_drilldown_" + current_id + "\" name=\"button_vm_drilldown\"><a href=\"/gcp/\">" + vm_map["vm_name"] + "</a></li>"

				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vm_name"] + "\" name=\"vm_name\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + "gcp" + "\" name=\"platform_type\">"
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["vpc_id"] + "\" name=\"vpc_id\" value=\"default_vpc\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"                        
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["region"] + "\" name=\"region\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["project_id"] + "\" name=\"project_id\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["fqdn"] + "\" name=\"fqdn\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["userid"] + "\" name=\"userid\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_ip_address"] + "\" name=\"public_ip_address\">"            
				html_row += "<input form=\"form_" + current_id + "\" type=\"hidden\" value=\"" + vm_map["private_ip_address"] + "\" name=\"private_ip_address\">"                         

				html_row += "</form>\n"

				$('#gcp_imported_vm_list').append(html_row);            

				$("#button_vm_drilldown_" +  current_id).click(function(event) {
					event.preventDefault();					
					var form_id = $(this).attr("form");
					document.getElementById(form_id).target = "operations_frame"; 
					document.getElementById(form_id).action = "/specific_vm/";
					document.getElementById(form_id).submit();
				});

				current_id += 1
			}); 
		},
  });
  /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

$(document).ready(function() {	

	get_nav_frame_maps();

});
