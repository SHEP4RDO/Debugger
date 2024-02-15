# Module mklog

The `mklog` module provides a toolkit for logging in Go applications with the ability to configure various aspects of log output and formatting.

## Key Features:

-   **Configuration Flexibility:** The module offers a configuration mechanism that allows tuning logging parameters such as debug mode, time format, console output, file output, and more.
    
-   **Log Format Choice:** Various log formats are supported, including Plain Text, JSON, XML, YAML, as well as the ability to define a custom format.
    
-   **Ease of Use:** A simple interface and the ability to initialize the logger via a configuration file simplify the process of integrating logging into an application.
    
-   **Logging Levels Support:** The module provides logging levels such as Debug, Trace, Info, Warning, Error, and Fatal, along with the ability to create custom levels, allowing detailed control over which messages are recorded in the log.
    
-   **Error Stack:** The ability to log detailed error information, including call stacks, when using corresponding methods.
    
-   **File-Based Configuration:** Easy setup of logging parameters through a configuration file in YAML, JSON, and XML formats, or by using a custom configuration.
    
This module ensures a reliable and flexible logging mechanism for your application, helping you efficiently monitor and analyze its operation.

## Possible Configuration Settings for the `mklog` Module include:

The module allows you to:
-   **Console Output:** Enable or disable log output to the console.
- -   **Debug Mode:** Enable or disable debug mode.
- -  **Date Format:** Set the date format for log messages.
 - -   **Detailed Error Output:** Enable or disable detailed error output, which includes error details in the stack.
  
- -   **Log Formatter Configuration:** Set the log formatter for the logger instance.
    
- -  **User-Defined Formatter:** Set a user-defined log formatter, providing users the ability to define their own format.
    
- -  **Set User-Defined Log Level Names:** Set custom log level names as provided by the user.

 -  **File Logging:** Enable or disable logging to a file and set the file path and name for recording.
-  -  **Log with Date in File Names:** Enable or disable adding the date to log file names.
- -  **Date Format in File Names:** Set the date format used in log file names.
    
- -  **Set Storage Path for Files:** Set the default path for log files, either using the default or explicitly specifying a path.
    
-  - **Set Log File Type:** Specify the file type (extension) for logs.

## Usage

Instructions on how to use the `mklog` module can be found in the corresponding [wiki](https://github.com/SHEP4RDO/mklog/wiki) of the project.
