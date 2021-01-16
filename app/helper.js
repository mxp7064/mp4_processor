// helper functions

const FileType = require('file-type');
const readChunk = require('read-chunk');
const { AppError, InvalidFileTypeError } = require('./errors.js')
const NATS = require('nats');

const isFileType = (filePath, requiredFileType) => {
  return new Promise(async (resolve, reject) => {
    let buffer;
    try {
      buffer = readChunk.sync(filePath, 0, 200);
      fileType = await FileType.fromBuffer(buffer);
    } catch (err) {
      reject(new AppError("Error occurred during file reading", err))
      return
    }

    if (fileType == undefined || fileType.mime != requiredFileType) {
      resolve(false)
      return
    }
    resolve(true)
  })
}

const requestOne = (natsClient, subject, message, timeout) => {
  return new Promise((resolve, reject) => {
    natsClient.requestOne(subject, message, timeout, response => {
      if (response instanceof NATS.NatsError) {
        reject(response)
        return
      }
      resolve(response);
    })
  })
}

const promptCLI = (commandLine, question) => {
  return new Promise((resolve, reject) => {
    commandLine.question(question, answer => {
      resolve(answer);
    })
  });
};

exports.isFileType = isFileType
exports.requestOne = requestOne
exports.promptCLI = promptCLI