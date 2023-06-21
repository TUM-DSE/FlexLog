#include <stdint.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif
    typedef void* PMLog;
	PMLog startUp();
	PMLog startUp_idx(int);
	void finalize(PMLog log);
    void PMLogFree(PMLog log);
    int cAppend(PMLog log, const char* record, uint64_t lsn);
    int cCommit(PMLog log,uint64_t lsn, uint64_t gsn) ;
    const char *cRead(PMLog log, uint64_t gsn, void *next_gsn);
    int cTrim(PMLog log, uint64_t gsn);
#ifdef __cplusplus
}
#endif
