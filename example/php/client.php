<?php
function go_pack($string){
	    //$string = json_encode($string);
	    $header = json_encode([
	        "content-length" => strlen($string),
	        "is-alive" => 0,
	    ]);
		return "[HDR]" . $header. "[/HDR]". $string;
}

$addr = "127.0.0.1";
$port = 8821;
$sock = @fsockopen($addr, $port, $errno, $errstr, 100);

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
