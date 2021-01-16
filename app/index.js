const readline = require('readline');
const NATS = require('nats');
const { requestOne, isFileType, promptCLI } = require('./helper.js')
const { NATS_URL, NATS_USER, NATS_PASS, NATS_SUBJECT, NATS_REQUEST_TIMEOUT } = require('./constants.js');
const natsClient = NATS.connect({ url: NATS_URL, user: NATS_USER, pass: NATS_PASS })
const commandLine = readline.createInterface(process.stdin, process.stdout);
const requiredFileType = "video/mp4";

(async () => {
  while (true) {
    let filePath = await promptCLI(commandLine, "Provide mp4 file path or type 'exit' > ")

    if (filePath == 'exit') {
      break;
    }

    try {
      let fileTypeCorrect = await isFileType(filePath, requiredFileType)
      if (!fileTypeCorrect) {
        console.log("File type must be " + requiredFileType)
        continue
      }
      let response = await requestOne(natsClient, NATS_SUBJECT, filePath, NATS_REQUEST_TIMEOUT)
      console.log(response)
    } catch (err) {
      // we can also log error stack from AppError's cause with err.stack
      console.error(err.message)
    }
  }

  natsClient.close();
  commandLine.close();
  console.log('Bye bye!');
})();

natsClient.on('error', (err) => {
  console.error('Error [' + natsClient.currentServer + ']:', err.message)
  process.exit(1)
})
