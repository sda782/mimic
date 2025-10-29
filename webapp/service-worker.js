const CACHE_NAME = "aion-cache-v1";
const urlsToCache = [
  "./",
  "./index.html",
  "./style.css",
  "./normalize.css",
  "./main.js",
  "./setup.html",
  "./assets/icons/app_icon_dark.svg"
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(urlsToCache))
  );
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => response || fetch(event.request))
  );
});