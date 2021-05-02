package config

const HttpPort = ":80"
const Environment = "env"
const GoogleClientKey = "599717309315-c84f5ijm874mu2of1i1g6qm6ufbfvmn4.apps.googleusercontent.com"
const GoogleSecret = "x9XbDukgssGemHHeni_UBckZ"
const GoogleAuthCallbackUrl = "http://localhost:5000/auth/google/callback?provider=google"
const FacebookClientKey = "1262553130621343"
const FacebookSecret = "330deb872bc3dec1438fe30feb74c766" //"edaf03f3aee9dc651f68e3ec50077a88"
const FacebookAuthCallbackUrl = "http://127.0.0.1/auth/facebook/callback?provider=facebook"

const DefaultMail = "mail2193"

//const DefaultMail = "mr.uwaifo@gmail.com"

const DefaultMailPassword = "JbK9mSraq8"

var DatabaseName = "testing"
var ConfirmationDbName = "confirmations"

var RedisDB = 0

var RedirectUrl = "http://localhost:8000"
