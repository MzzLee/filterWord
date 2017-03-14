<?php
	$sock = @fsockopen("127.0.0.1", "9901", $errno,
            $errstr, 100);
	
	function go_pack($string){
	    $string = json_encode($string);
		return "[Header]".strlen($string)."[/Header]".$string;
	}
	$body = go_pack($argv[1]);
	fwrite($sock, $body, strlen($body));
	$response = fgets($sock, 1024);
	@socket_close($sock);
	print_r($response."\r\n");
?>
