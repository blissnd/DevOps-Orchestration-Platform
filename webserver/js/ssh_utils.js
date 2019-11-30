$(document).ready(function() {	
	
	$('#ssh_window_button_close').click(function() {	  
	  var show_ssh = "false";
	  $('#ssh_frame').remove();
	  $('#ssh_window').hide();
	});	
	
	$('#ssh_window_button_detach').click(function() {	  
	  
	  var ssh_frame_handle = document.createElement('iframe');
	  ssh_frame_handle.setAttribute('name', 'ssh_popup_frame');
	  ssh_frame_handle.setAttribute('id', 'ssh_frame');
	  ssh_frame_handle.setAttribute('style', 'height:300px;float:left;display:inline-block;frameborder:5;scrolling:yes;width:100%;resize:both;');
	  ssh_frame_handle.setAttribute('src', '/ssh_popup_frame');
	
		var new_ssh_window = ssh_frame_handle.src;
    window.open(new_ssh_window, selected_vm + "-" + cloud_platform + "_" + instance_construction_type + "_" + docker_container_name, 'width=800,height=500,left=0,top=100,screenX=0,screenY=100')	  
    
    new_ssh_window.selected_vm = selected_vm
    new_ssh_window.cloud_platform = cloud_platform
    new_ssh_window.instance_construction_type = instance_construction_type

    var show_ssh = "false";
	  $('#ssh_frame').remove();
	  $('#ssh_window').hide();
	});
	
});
