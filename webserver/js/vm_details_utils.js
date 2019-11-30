function get_vm_list_for_all_platforms() {

	$.ajax({
		url: "/get_aws_imported_map/",
		type: 'post',
		dataType: 'json',
		
		success : function(vm_list_1) {
			$.ajax({
				url: "/get_aws_map/",
				type: 'post',
				dataType: 'json',
	
				success : function(vm_list_2) {

					$.ajax({
						url: "/get_gcp_imported_map/",
						type: 'post',
						dataType: 'json',
						
						success : function(vm_list_3) {

							$.ajax({
								url: "/get_gcp_map/",
								type: 'post',
								dataType: 'json',
								
								success : function(vm_list_4) {

									$.ajax({
										url: "/get_vm_map/",
										type: 'post',
										dataType: 'json',
										
										success : function(vm_list_5) {

											var client_selection_block = document.getElementById("server_form_client_server_select_block");
											var server_selection_block = document.getElementById("client_form_client_server_select_block");

											$.each(vm_list_1, function(vm_name, vm_map) {
												client_selection_block.options[client_selection_block.options.length] = new Option(vm_name, vm_name);
												server_selection_block.options[server_selection_block.options.length] = new Option(vm_name, vm_name);
											});
											$.each(vm_list_2, function(vm_name, vm_map) {
												client_selection_block.options[client_selection_block.options.length] = new Option(vm_name, vm_name);
												server_selection_block.options[server_selection_block.options.length] = new Option(vm_name, vm_name);
											});
											$.each(vm_list_3, function(vm_name, vm_map) {
												client_selection_block.options[client_selection_block.options.length] = new Option(vm_name, vm_name);
												server_selection_block.options[server_selection_block.options.length] = new Option(vm_name, vm_name);
											});
											$.each(vm_list_4, function(vm_name, vm_map) {
												client_selection_block.options[client_selection_block.options.length] = new Option(vm_name, vm_name);
												server_selection_block.options[server_selection_block.options.length] = new Option(vm_name, vm_name);
											});
											$.each(vm_list_5, function(vm_name, vm_map) {
												client_selection_block.options[client_selection_block.options.length] = new Option(vm_name, vm_name);
												server_selection_block.options[server_selection_block.options.length] = new Option(vm_name, vm_name);
											});
										},
									});
								},
							});
						},
					});
				},
			});
		},
	});
}
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
function get_platform_specific_html(section_name, html_row, vm_name, current_id, ajax_url_string) {

	$.ajax({
		url: ajax_url_string,
		type: 'post',
		dataType: 'json',
		
		success : function(data) {

			vm_map = data[vm_name];			

			html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + vm_name + "\" name=\"vm_name\">\n"
			html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + vm_map["instance_construction_type"] + "\" name=\"instance_construction_type\">"
      html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + vm_map["platform_type"] + "\" name=\"platform_type\">\n"
      html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + docker_container_name + "\" name=\"docker_container_name\">\n"      
      html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + vm_map["fqdn"] + "\" name=\"fqdn\">"            
      html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + vm_map["public_ip_address"] + "\" name=\"public_ip_address\">"
      html_row += "<input form=\"" + current_id + "\" type=\"hidden\" value=\"" + vm_map["OS"] + "\" name=\"OS\">"      

			html_row += "<input form=\"" + current_id + "\" type=\"hidden\" id=\"action_to_execute_" + current_id + "\" name=\"ansible_deploy\" value=\"Deploy\">"
			html_row += "<input form=\"" + current_id + "\" type=\"submit\" name=\"generate_inventory\" id=\"generate_inventory_" + current_id + "\" value=\"Generate Inventory\">&nbsp;&nbsp;"
      html_row += "<input form=\"" + current_id + "\" type=\"submit\" name=\"ansible_deploy_button\" id=\"ansible_deploy_" + current_id + "\" value=\"Deploy\">"
			html_row += "</form>\n"

			$('#' + section_name).append(html_row);

			////////////////////////////////////////////////////////////////////////////
			
			$("#ansible_deploy_" +  current_id).click(function(event) {		

				var form_id = $(this).attr("form");				

				event.preventDefault();								

				document.getElementById("action_to_execute" + "_" + form_id).setAttribute("name", "ansible_deploy");
				document.getElementById(form_id).action = "/ansible/run_ansible/";
				document.getElementById(form_id).target = "_blank";
				
				if ($('#' + form_id)[0].checkValidity()) {
					document.getElementById(form_id).submit();
				}

			});

			$("#generate_inventory_" +  current_id).click(function(event) {				

				var form_id = $(this).attr("form");

				$('#' + form_id)[0].checkValidity()

				event.preventDefault();				

				document.getElementById("action_to_execute" + "_" + form_id).setAttribute("name", "generate_inventory");
				document.getElementById(form_id).action = "/specific_vm/";
				document.getElementById(form_id).target = "";												 								
				
				if ($('#' + form_id)[0].checkValidity()) {
					document.getElementById(form_id).submit();
				}

			});

			if (list_populate_counter == 0) {
				get_vm_list_for_all_platforms();
				list_populate_counter++;
			}
		},
	});
}
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
function get_map_for_correct_cloud_platform(section_name, html_row, cloud_platform, selected_vm, current_id, instance_construction_type) {	

	var ajax_url_string;

	if (instance_construction_type == "Imported" && cloud_platform == "aws") {
		ajax_url_string = "/get_aws_imported_map/";
	} else if (instance_construction_type == "Managed" && cloud_platform == "aws") {
		ajax_url_string = "/get_aws_map/";
	} else if (instance_construction_type == "Imported" && cloud_platform == "gcp") {
		ajax_url_string = "/get_gcp_imported_map/";
	} else if (instance_construction_type == "Managed" && cloud_platform == "gcp") {
		ajax_url_string = "/get_gcp_map/";
	} else {
		ajax_url_string = "/get_vm_map/"
	}

	get_platform_specific_html(section_name, html_row, selected_vm, current_id, ajax_url_string)		
}
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
$(document).ready(function() {

	list_populate_counter = 0;

	$('#docker_form').show();
	$('#iptables_form').show();

	if (docker_container_display != "") {
		$('#docker_form').hide();
		$('#iptables_form').hide();
	}

	//////////////////////////////////////////////////// BASIC PLAYBOOKS /////////////////////////////////////////////////////

	var current_form_id = "basic_playbook_form";

	var html_row = "<form id=\"" + current_form_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"

	html_row += "<select id=\"" + current_form_id + "_playbook_select_block\" form=\"" + current_form_id + "\" name=\"deployment_list_selection\" required>"	
	html_row += "<option  value=\"\">&lt;Select Playbook&gt;</option>"

	html_row += "<option  value=\"deployment_list_basic_config\">Deploy Basic Config</option>"
	html_row += "<option  value=\"deployment_list_net_tools\">Deploy net-tools</option>"
	html_row += "<option  value=\"deployment_list_set_aws_hostname\">Set Host Name</option>"
	html_row += "<option  value=\"deployment_list_deploy_iptables\">Deploy Firewall</option>"
	html_row += "<option  value=\"deployment_list_mongodb\">Deploy MongoDB</option>"
	html_row += "<option  value=\"deployment_list_java\">Install Java</option>"

	html_row += "</select>"

	get_map_for_correct_cloud_platform("ansible_deployment_form_basic", html_row, cloud_platform, selected_vm, current_form_id, instance_construction_type);		

	///////////////////////////////////////////////////////// SERVERS ///////////////////////////////////////////////////////

	var current_form_id = "server_form";

	html_row2 = "<form id=\"" + current_form_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"

	html_row2 += "<select id=\"" + current_form_id + "_playbook_select_block\" form=\"" + current_form_id + "\" name=\"deployment_list_selection\" required>"
	html_row2 += "<option  value=\"\">&lt;Select Playbook&gt;</option>"

	html_row2 += "<option  value=\"deployment_list_DNS\">Deploy DNS Server</option>"
	html_row2 += "<option  value=\"deployment_list_LDAP\">Deploy LDAP Server</option>"
	html_row2 += "<option  value=\"deployment_list_nexus\">Deploy Sonatype Nexus 3</option>"
	html_row2 += "<option  value=\"deployment_list_prometheus\">Deploy Prometheus</option>"
	html_row2 += "<option  value=\"deployment_list_grafana\">Deploy Grafana</option>"
	html_row2 += "<option  value=\"deployment_list_x11\">Deploy X11 & VNC</option>"
	html_row2 += "<option  value=\"deployment_list_logstash\">Deploy LogStash</option>"
	html_row2 += "<option  value=\"deployment_list_elasticsearch\">Deploy ElasticSearch</option>"
	html_row2 += "<option  value=\"deployment_list_kibana\">Deploy Kibana</option>"
	html_row2 += "<option  value=\"deployment_list_kubernetes_preparation\">Kubernetes Preparation</option>"
	html_row2 += "<option  value=\"deployment_list_kubernetes_master\">Kubernetes Master</option>"

	html_row2 += "</select>"

	html_row2 += "<select multiple id=\"" + current_form_id + "_client_server_select_block\" form=\"" + current_form_id + "\" name=\"selection_list_of_clients\" required>"	

	get_map_for_correct_cloud_platform("ansible_deployment_form_servers", html_row2, cloud_platform, selected_vm, current_form_id, instance_construction_type);

	////////////////////////////////////////////////////////// CLIENTS ////////////////////////////////////////////////////////

	var current_form_id = "client_form";

	html_row3 = "<form id=\"" + current_form_id + "\" action=\"/\" target=\"_blank\"  method=\"POST\">\n"

	html_row3 += "<select id=\"" + current_form_id + "_playbook_select_block\" form=\"" + current_form_id + "\" name=\"deployment_list_selection\" required>"
	html_row3 += "<option  value=\"\">&lt;Select Playbook&gt;</option>"

	html_row3 += "<option  value=\"deployment_list_DNS_Client\">Deploy DNS Client</option>"
	html_row3 += "<option  value=\"deployment_list_LDAP_Client\">Deploy LDAP Client</option>"
	html_row3 += "<option  value=\"deployment_list_docker_install\">Install Docker</option>"
	html_row3 += "<option  value=\"deployment_list_configure_docker_registry_client\">Configure Docker Client</option>"
	html_row3 += "<option  value=\"deployment_list_node_exporter\">Deploy node-exporter</option>"
	html_row3 += "<option  value=\"deployment_list_filebeat\">Deploy Filebeat</option>"
	html_row3 += "<option  value=\"deployment_list_kubernetes_agent\">Kubernetes Agent</option>"

	html_row3 += "</select>"

	html_row3 += "<select id=\"" + current_form_id + "_client_server_select_block\" form=\"" + current_form_id + "\" name=\"server_selection\" required>"	
	html_row3 += "<option  value=\"\">&lt;Select Server for Client&gt;</option></select>"

	get_map_for_correct_cloud_platform("ansible_deployment_form_clients", html_row3, cloud_platform, selected_vm, current_form_id, instance_construction_type);

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	//Hide ssh window frame
	if (show_ssh == "false") {
	  
	  $('#ssh_window').hide();

	} else {

	  $('#ssh_window').show();

	  if (selected_vm != "") {

	    $('#ssh_div_container').html("<iframe name=\"ssh_frame\" id=\"ssh_frame\" style=\"height:300px;float:left;display:inline-block;frameborder:5;scrolling:yes;width:100%;resize:both;\" src=\"/ssh_frame\"></iframe>");
	  }
	}

	$("#button_sshlaunch_0").click(function() {		

    var form_id = $(this).attr("form");
    document.getElementById(form_id).action = "/specific_vm/";
    document.getElementById(form_id).target = ""; 
    document.getElementById(form_id).submit();
  });

	$("#button_iptables_1").click(function() {

    var form_id = $(this).attr("form");
    document.getElementById(form_id).action = "/iptables_security/";
    document.getElementById(form_id).target = ""; 
    document.getElementById(form_id).submit();
  });
  
  $("#button_docker_2").click(function() {

    var form_id = $(this).attr("form");
    document.getElementById(form_id).action = "/docker/";
    document.getElementById(form_id).target = ""; 
    document.getElementById(form_id).submit();
  });

  $("#button_get_ports_3").click(function(event) {

  	event.preventDefault();

    var form_id = $(this).attr("form");

    var form_get_ports_vm_name = $("#form_get_ports_vm_name").val()
    var form_get_ports_platform_type = $("#form_get_ports_platform_type").val()
    var form_get_ports_instance_construction_type = $("#form_get_ports_instance_construction_type").val()
    var form_get_ports_docker_container_name = $("#form_get_ports_docker_container_name").val()			  
	  
	  webform_obj = {	vm_name: form_get_ports_vm_name,
                        platform_type: form_get_ports_platform_type,
                        instance_construction_type: form_get_ports_instance_construction_type,
                        docker_container_name: form_get_ports_docker_container_name
                     	};   

	  var listening_port_frame_handle = document.createElement('iframe');
		listening_port_frame_handle.setAttribute('name', 'listening_port_frame_handle');
		listening_port_frame_handle.setAttribute('id', 'listening_port_frame_handle');
		listening_port_frame_handle.setAttribute('style', 'height:300px;float:left;display:inline-block;frameborder:5;scrolling:yes;width:100%;resize:both;');
		listening_port_frame_handle.setAttribute('src', '/listening_port_popup_frame/');

  	var port_frame_window = listening_port_frame_handle.src;	    
  	
  	window.open(port_frame_window, selected_vm + "-" + cloud_platform + "_" + instance_construction_type + "_" + docker_container_name, 'width=400,height=400,left=0,top=100,screenX=0,screenY=100')  	
  	port_frame_window.webform_obj = webform_obj
	  
  });
  ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
});
