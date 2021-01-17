const readline = require('readline');
const NATS = require('nats');
const { requestOne, isFileType, promptCLI } = require('./helper.js')
const { NATS_URL, NATS_USER, NATS_PASS, NATS_SUBJECT, NATS_REQUEST_TIMEOUT } = require('./constants.js');
const natsClient = NATS.connect({ url: NATS_URL, user: NATS_USER, pass: NATS_PASS })
const commandLine = readline.createInterface(process.stdin, process.stdout);
const requiredFileType = "video/mp4";

/*
  main program function which indefinitely waits for user input on the CLI
  it validates the user input which must be a valid mp4 file path
  sends the file path via nats on the appropriate subject and prints the response
 */
(async () => {
  while (true) {
    let filePath = await promptCLI(commandLine, "Provide mp4 file path or type 'exit' > ")

    if (filePath == "") {
      console.log("Please provide mp4 file path")
      continue
    }

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

  // when we break the above loop, we close nats and the command line and effectively exit the program
  natsClient.close();
  commandLine.close();
  console.log('Bye bye!');
})();

// nats error handling
natsClient.on('error', (err) => {
  console.error('Error [' + natsClient.currentServer + ']:', err.message)
  process.exit(1)
})
