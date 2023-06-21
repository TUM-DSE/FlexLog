#include "PMLog.hpp"
#include "PMLog.h"

void PMLogFree(PMLog log) {
	cppPMLog *cppLog = reinterpret_cast<cppPMLog *>(log);

	delete cppLog;
}

PMLog startUp() {
	return static_cast<PMLog>(cppStartUp());
}

PMLog startUp_idx(int idx) {
	return static_cast<PMLog>(cppStartUp(idx));
}

void finalize(PMLog log) {
	persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

	cppFinalize(cppLog);
}

int cAppend(PMLog log, const char* record, uint64_t lsn) {
	persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

	return cppLog->Append(record, lsn);
}

int cCommit(PMLog log, uint64_t lsn, uint64_t gsn){
	persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

	return cppLog->Commit(lsn, gsn);
}

const char *cRead(PMLog log, uint64_t gsn, void* next_gsn) {
	persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);
	uint64_t *nxtgsn = reinterpret_cast<uint64_t *>(next_gsn);

	return cppLog->Read(gsn, nxtgsn);
}

int cTrim(PMLog log, uint64_t gsn) {
	persistent_ptr<cppPMLog> cppLog = persistent_ptr<cppPMLog>((pmem::detail::sp_element<cppPMLog>::type *) log);

	return cppLog->Trim(gsn);
}
