// const fs = require('fs');

// function convertCookies(jsonFilePath, netscapeFilePath) {
//     fs.readFile(jsonFilePath, 'utf8', (err, data) => {
//         if (err) {
//             console.error("Error reading JSON file:", err);
//             return;
//         }

//         const cookies = JSON.parse(data);
//         const netscapeCookies = [];

//         // Add Netscape header
//         netscapeCookies.push("# Netscape HTTP Cookie File");
//         netscapeCookies.push("# This is a generated file! Do not edit.");
//         netscapeCookies.push("# http://www.netscape.com/newsref/std/cookie_spec.html");

//         cookies.forEach(cookie => {
//             const domain = cookie.domain.startsWith(".") ? cookie.domain : `.${cookie.domain}`;
//             const includeSubdomains = "TRUE"; // Always TRUE for domain cookies
//             const path = cookie.path || "/";
//             const secure = cookie.secure ? "TRUE" : "FALSE";
//             const expiration = cookie.expirationDate ? Math.floor(cookie.expirationDate) : "0"; // Unix timestamp
//             const name = cookie.name;
//             const value = cookie.value;

//             const cookieLine = [
//                 domain,
//                 includeSubdomains,
//                 path,
//                 secure,
//                 expiration,
//                 name,
//                 value
//             ].join("\t");

//             netscapeCookies.push(cookieLine);
//         });

//         // Write the Netscape cookies to the file
//         const cookiesContent = netscapeCookies.join("\n");
//         fs.writeFile(netscapeFilePath, cookiesContent, 'utf8', (err) => {
//             if (err) {
//                 console.error("Error writing Netscape file:", err);
//             } else {
//                 console.log(`Cookies converted to Netscape format and saved to ${netscapeFilePath}`);
//             }
//         });
//     });
// }

// convertCookies("cookies.json", "cookies.txt");