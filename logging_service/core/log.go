package log

// Writing logs to files:

/* 
 * 1. Get payload
 * 2. Check severity level
 * 3. Check message counts
 * 4. If message count is greater equal to LOG_FILE_MAX_COUNT for the current log file, create a new file.
 * 5. Otherwise, use the most recent log file for that severity level.
 * 6. Injest log message, write to file, close file.
 */

// Reading log files

/* How you could query:
 * 	date range
 * 	log number/id
 * 	severity level
 * 	message content
 * 	log types
 */

/*
 * 1. Get payload
 * 2. Check for log number/id
 * 		If present, get that specific log.
 *		otherwise, continue.
 * 3. Check for date range
 *		If present, add to search query.
 * 4. Check for severity level(s).
 * 5. If hitting the /logs endpoint instead of /logs/error or /logs/debug, check for log type(s).
 * 6. Check severity level. If present, search only the folder(s) with the severity level(s) given/.
 * 	  Otherwise, check all severity levels.
 * 7. If message content is present, filter out after searching all of these. (Probably will not be implimented).
 */
