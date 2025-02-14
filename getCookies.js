const fs = require('fs');

// Funzione per convertire i cookies da JSON a Netscape
function convertCookiesToNetscape(cookies) {
  let netscapeCookies = [
    "# Netscape HTTP Cookie File\n",
    "# This file was generated from JSON\n",
    "# https://www.netscape.com\n"
  ];

  cookies.forEach(cookie => {
    let cookieLine = [
      ".tiktok.com", 
      cookie.name,   
      cookie.value,  
      "",            
      "FALSE",       
      "FALSE",       
      "0",          
      cookie.expires || "0" 
    ].join("\t");

    netscapeCookies.push(cookieLine);
  });

  return netscapeCookies.join("\n");
}

fs.readFile('cookies.json', 'utf8', (err, data) => {
  if (err) {
    console.error("Errore nella lettura del file JSON:", err);
    return;
  }

  const cookies = JSON.parse(data);

  const netscapeCookies = convertCookiesToNetscape(cookies);

  fs.writeFile('cookies.txt', netscapeCookies, (err) => {
    if (err) {
      console.error("Errore nella scrittura del file Netscape:", err);
    } else {
      console.log('File di cookies salvato come cookies.txt');
    }
  });
});
