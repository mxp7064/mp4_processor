// custom Error type for our application
class AppError extends Error {
  constructor(message, cause) {
    super(message + (cause ? ": " + cause.message : ""))
    this.stack = cause?.stack
  }
}

exports.AppError = AppError