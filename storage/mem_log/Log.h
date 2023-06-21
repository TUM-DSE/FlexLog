#include <stdint.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif
    typedef void* Log;
    Log LogNew();
    void LogFree(Log log);
    int cAppend(Log log, const char* record, uint64_t lsn);
    int cCommit(Log log,uint64_t lsn, uint64_t gsn) ;
    const char *cRead(Log log, uint64_t gsn, void *next_gsn);
    int cTrim(Log log, uint64_t gsn);
#ifdef __cplusplus
}
#endif