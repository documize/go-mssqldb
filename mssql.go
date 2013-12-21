package mssql

import (
    "database/sql"
    "database/sql/driver"
//	"math"
//	"math/big"
//	"time"
//	"unsafe"
    "strings"
)

func init() {
    sql.Register("go-mssql", &MssqlDriver{})
}

type MssqlDriver struct {
}

type MssqlConn struct {
    sess *TdsSession
}

type MssqlTx struct {
	c *MssqlConn
}

//func (tx *MssqlTx) Commit() error {
//	_, err := oleutil.CallMethod(tx.c.db, "CommitTrans")
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (tx *MssqlTx) Rollback() error {
//	_, err := oleutil.CallMethod(tx.c.db, "Rollback")
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (c *MssqlConn) exec(cmd string) error {
//	_, err := oleutil.CallMethod(c.db, "Execute", cmd)
//	return err
//}
//
func (c *MssqlConn) Begin() (driver.Tx, error) {
//	_, err := oleutil.CallMethod(c.db, "BeginTrans")
//	if err != nil {
//		return nil, err
//	}
//	return &AdodbTx{c}, nil
    return nil, nil
}

func parseConnectionString(dsn string) (res map[string]string) {
    res = map[string]string{}
    parts := strings.Split(dsn, ";")
    for _, part := range parts {
        if len(part) == 0 {
            continue
        }
        lst := strings.SplitN(part, "=", 2)
        name := strings.ToLower(lst[0])
        if len(name) == 0 {
            continue
        }
        var value string = ""
        if len(lst) > 1 {
            value = lst[1]
        }
        res[name] = value
    }
    return res
}

func (d *MssqlDriver) Open(dsn string) (driver.Conn, error) {
    params := parseConnectionString(dsn)
    buf, err := Connect(params)
    if err != nil {
        return nil, err
    }
    return &MssqlConn{buf}, nil
}

func (c *MssqlConn) Close() error {
    return c.sess.buf.transport.Close()
}

type MssqlStmt struct {
    c  *MssqlConn
    query string
}

func (c *MssqlConn) Prepare(query string) (driver.Stmt, error) {
	return &MssqlStmt{c, query}, nil
}

//func (s *MssqlStmt) Bind(bind []string) error {
//	s.b = bind
//	return nil
//}

func (s *MssqlStmt) Close() error {
//	s.s.Release()
    return nil
}

func (s *MssqlStmt) NumInput() int {
//	if s.b != nil {
//		return len(s.b)
//	}
//	_, err := oleutil.CallMethod(s.ps, "Refresh")
//	if err != nil {
//		return -1
//	}
//	val, err := oleutil.GetProperty(s.ps, "Count")
//	if err != nil {
//		return -1
//	}
//	c := int(val.Val)
//	return c
    return 0
}

//func (s *MssqlStmt) bind(args []driver.Value) error {
//	if s.b != nil {
//		for i, v := range args {
//			var b string = "?"
//			if len(s.b) < i {
//				b = s.b[i]
//			}
//			unknown, err := oleutil.CallMethod(s.s, "CreateParameter", b, 12, 1)
//			if err != nil {
//				return err
//			}
//			param := unknown.ToIDispatch()
//			defer param.Release()
//			_, err = oleutil.PutProperty(param, "Value", v)
//			if err != nil {
//				return err
//			}
//			_, err = oleutil.CallMethod(s.ps, "Append", param)
//			if err != nil {
//				return err
//			}
//		}
//	} else {
//		for i, v := range args {
//			var varval ole.VARIANT
//			varval.VT = ole.VT_I4
//			varval.Val = int64(i)
//			val, err := oleutil.CallMethod(s.ps, "Item", &varval)
//			if err != nil {
//				return err
//			}
//			item := val.ToIDispatch()
//			defer item.Release()
//			_, err = oleutil.PutProperty(item, "Value", v)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}

func (s *MssqlStmt) Query(args []driver.Value) (driver.Rows, error) {
//	if err := s.bind(args); err != nil {
//		return nil, err
//	}
    headers := []headerStruct{
        {hdrtype: dataStmHdrTransDescr,
         data: transDescrHdr{0, 1}.pack()},
    }
    if err := sendSqlBatch72(s.c.sess.buf, s.query, headers); err != nil {
        return nil, err
    }
    return &MssqlRows{sess: s.c.sess}, nil
}

func (s *MssqlStmt) Exec(args []driver.Value) (driver.Result, error) {
//	if err := s.bind(args); err != nil {
//		return nil, err
//	}
//	_, err := oleutil.CallMethod(s.s, "Execute")
//	if err != nil {
//		return nil, err
//	}
    return driver.ResultNoRows, nil
}

type MssqlRows struct {
    sess *TdsSession
//	s    *AdodbStmt
//	rc   *ole.IDispatch
    nc   int
    cols []string
}

func (rc *MssqlRows) Close() error {
//	_, err := oleutil.CallMethod(rc.rc, "Close")
//	if err != nil {
//		return err
//	}
    return nil
}

func (rc *MssqlRows) Columns() (res []string) {
    if !rc.sess.gotColumns {
        err := processResponseGetHeader(rc.sess)
        if err != nil {
            return []string{}
        }
    }
    res = make([]string, len(rc.sess.columns))
    for i := range rc.sess.columns {
        res[i] = rc.sess.columns[i].ColName
    }
    return res
}

func (rc *MssqlRows) Next(dest []driver.Value) (err error) {
    err = getRow(rc.sess)
    if err != nil {
        return err
    }
    for i := range dest {
        dest[i] = rc.sess.lastRow[i]
    }
    return nil
}
