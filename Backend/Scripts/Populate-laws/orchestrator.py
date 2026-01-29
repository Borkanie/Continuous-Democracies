import sys
import os
import time
import datetime

SCRIPT_DIR = os.path.dirname(__file__)
if SCRIPT_DIR not in sys.path:
	sys.path.insert(0, SCRIPT_DIR)

from log_writer import file_log

# Direct imports (files live in this directory)
from import_politicians import importPoliticians
from import_law import goFromLastforward
from activate_politicians import activate_all_groups
from dataBase_interaction import init_db, close_db

def run_cycle():
	year = datetime.datetime.now().year
	file_log(f"Orchestrator cycle start: {datetime.datetime.now().isoformat()}", alsoPrint=True)

	# 1) Import politicians only for the current year (IDs 1..333)
	try:
		file_log(f"Importing politicians for year {year}", alsoPrint=True)
		importPoliticians(1, 333, year)
	except Exception as e:
		file_log(f"Failed to import politicians: {e}", alsoPrint=True)

	# 2) Import new laws from last forward, attempt up to 100
	try:
		file_log("Calling import_law.goFromLastforward(100)", alsoPrint=True)
		goFromLastforward(100)
	except Exception as e:
		file_log(f"Failed to import laws: {e}", alsoPrint=True)

	file_log(f"Orchestrator cycle end: {datetime.datetime.now().isoformat()}", alsoPrint=True)

def main():
	file_log("Starting orchestrator (runs once every hour).", alsoPrint=True)

	# initialize shared DB connection once
	try:
		init_db()
		file_log("Database connection initialized.", alsoPrint=True)
	except Exception as e:
		file_log(f"Failed to initialize DB connection: {e}", alsoPrint=True)

	# trigger activation immediately on first loop by setting last_activate in the past
	last_activate = datetime.datetime.now() - datetime.timedelta(hours=12)

	try:
		while True:
			run_cycle()

			# every 12 hours run the activate script
			now = datetime.datetime.now()
			if (now - last_activate) >= datetime.timedelta(hours=12):
				try:
					file_log("Running activate_politicians.activate_all_groups()", alsoPrint=True)
					activate_all_groups()
					last_activate = now
				except Exception as e:
					file_log(f"Failed to activate politicians: {e}", alsoPrint=True)

			time.sleep(3600)
	except KeyboardInterrupt:
		file_log("Orchestrator stopped by user.", alsoPrint=True)
	finally:
		try:
			close_db()
			file_log("Database connection closed.", alsoPrint=True)
		except Exception:
			pass

if __name__ == "__main__":
	main()

