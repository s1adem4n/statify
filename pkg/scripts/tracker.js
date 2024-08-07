const serverAddress = "{{ .serverAddress }}";
// generate a session ID, or get it from local storage
let sessionID = localStorage.getItem("sessionID");
if (!sessionID) {
  sessionID = crypto.randomUUID();
  localStorage.setItem("sessionID", sessionID);
}

// get the ip address of the client
let ip = "";
const ipPromise = fetch(`${serverAddress}/ip`).then((response) =>
  response.text().then((text) => (ip = text))
);

const domain = window.location.hostname;
const userAgent = navigator.userAgent;

const isMobileDevice = () => {
  const isMobileUserAgent =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      userAgent
    );
  const isMobiUserAgent = /mobi/i.test(userAgent.toLowerCase());

  const isTouchScreen =
    "ontouchstart" in window ||
    navigator.msMaxTouchPoints ||
    navigator.maxTouchPoints > 0;

  return (isMobileUserAgent || isMobiUserAgent) && isTouchScreen;
};

const deviceType = isMobileDevice() ? "mobile" : "desktop";

const trackView = async () => {
  await ipPromise;

  const path = window.location.pathname + window.location.search;

  const data = {
    domain,
    path,
    ip,
    userAgent,
    session: sessionID,
    device: deviceType,
  };

  const res = await fetch(`${serverAddress}/api/collections/views/records`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (res.status !== 200) {
    console.error("Failed to track view");
  }

  console.log("Tracked view", data);
};

trackView();

window.addEventListener("pushstate", trackView);
window.addEventListener("replacestate", trackView);
window.addEventListener("popstate", trackView);
