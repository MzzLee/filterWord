<?php
    $addr = "10.242.92.17";
    $port = 8821;
	$sock = @fsockopen($addr, $port, $errno, $errstr, 100);
	function go_pack($string){
	    $string = json_encode($string);
	    $header = json_encode([
	        "content-length" => strlen($string),
	        "is-alive" => false,
	    ]);
		return "[header]" . $header. "[/header]". $string;
	}

	$content = isset($argv[1]) ? $argv[1] : '';
	$body = go_pack($content);
	if($sock){
		fwrite($sock, $body, strlen($body));
    	$response = fgets($sock, 1024);
    	fclose($sock);
    	print_r($response."\r\n");
	}else{
	    print_r("Connecting Sock: {$errstr}!\r\n");
	}

?>
