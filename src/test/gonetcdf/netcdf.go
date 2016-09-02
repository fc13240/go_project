// Package GoNetCDF is an interface to the netCDF data format.
// For more information, refer to the netCDF C interface guide:
// http://www.unidata.ucar.edu/software/netcdf/docs/netcdf-c/
package gonetcdf

// #cgo LDFLAGS: -lnetcdf
// #include <stdlib.h>
// #include <netcdf.h>
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

const (
	// Dummy string for converting between C and go string formats.
	// If it isn't long enough, we end up with memory errors.
	cstr = "`````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````"
)

// This function creates a new netCDF dataset, returning a netCDf.ID
//that can subsequently be used to refer to the netCDF dataset in 
//other netCDF function calls. The new netCDF dataset opened for 
//write access and placed in define mode, ready for you to add 
//dimensions, variables, and attributes. 
func Create(fname string, mode string) (ncf *NCfile, err error) {
	Cfname := C.CString(fname)
	ncf = new(NCfile)
	defer C.free(unsafe.Pointer(Cfname))
	ncid := C.int(0)
	var cmode C.int
	switch mode {
	case "noclobber":
		cmode = C.int(C.NC_NOCLOBBER)
	case "share":
		cmode = C.int(C.NC_SHARE)
	case "64bitoffset":
		cmode = C.int(C.NC_64BIT_OFFSET)
	case "netcdf4":
		cmode = C.int(C.NC_NETCDF4)
	case "classicmodel":
		cmode = C.int(C.NC_CLASSIC_MODEL)
	default:
		err = errors.New("Incorrect file mode " + mode +
			". Options are noclobber, share, 64bitoffset, netcdf4 and classicmodel")
		return
	}
	e := C.nc_create(Cfname, cmode, &ncid)
	err = errAnal(e)
	if err != nil {
		return
	}
	ncf.Fname = fname
	ncf.ID = ncid
	ncf.VarNames = make(map[string]int)
	ncf.DimNames = make(map[string]int)
	ncf.GlobalAttNames = make(map[string]int)
	ncf.GlobalAttString = make(map[string]string)
	ncf.GlobalAttInt = make(map[string]int)
	ncf.GlobalAttFloat = make(map[string]float64)
	return
}

// The function nc_open opens an existing netCDF dataset for access. It determines the underlying file format automatically. Use the same call to open a netCDF classic, 64-bit offset, or netCDF-4 file. 
func Open(fname string, mode string) (ncf *NCfile, err error) {
	Cfname := C.CString(fname)
	ncf = new(NCfile)
	defer C.free(unsafe.Pointer(Cfname))
	ncid := C.int(0)
	var cmode C.int
	switch mode {
	case "nowrite":
		cmode = C.int(C.NC_NOWRITE)
	case "write":
		cmode = C.int(C.NC_WRITE)
	case "share":
		cmode = C.int(C.NC_SHARE)
	default:
		err = errors.New("Incorrect open mode " + mode +
			". Options are nowrite, write, and share")
		return
	}
	e := C.nc_open(Cfname, cmode, &ncid)
	err = errAnal(e)
	if err != nil {
		return
	}
	ncf.Fname = fname
	ncf.ID = ncid

	// get some information about the file
	ndims, nvars, ngatts, _, err := ncf.inq()
	if err != nil {
		return
	}
	// Get the names and ID numbers of all global attributes.
	ncf.DimNames = make(map[string]int)
	ncf.GlobalAttNames = make(map[string]int)
	ncf.GlobalAttString = make(map[string]string)
	ncf.GlobalAttInt = make(map[string]int)
	ncf.GlobalAttFloat = make(map[string]float64)
	varid := C.NC_GLOBAL
	for n := 0; n < ngatts; n++ {
		var name string
		name, err = ncf.inqAttName(varid, n)
		if err != nil {
			return
		}
		ncf.GlobalAttNames[name] = n
		var attType string
		attType, err = ncf.InqAttType(name, "global")
		if err != nil {
			return
		}
		switch attType {
		case "int":
			var val int
			val, err = ncf.GetAttInt(name, "global")
			ncf.GlobalAttInt[name] = val
			if err != nil {
				return
			}
		case "float":
			var val float64
			val, err = ncf.GetAttFloat(name, "global")
			ncf.GlobalAttFloat[name] = val
			if err != nil {
				return
			}
		case "string":
			var val string
			val, err = ncf.GetAttString(name, "global")
			ncf.GlobalAttString[name] = val
			if err != nil {
				return
			}
		default:
			err = errors.New("Can't handle attribute type " +
				attType)
			if err != nil {
				return
			}
		}
	}
	// Get IDs of all the dimensions
	for n := 0; n < ndims; n++ {
		var name string
		name, err = ncf.inqDimName(n)
		if err != nil {
			return
		}
		ncf.DimNames[name] = n
	}

	// Get the names and ID numbers of all variables.
	ncf.VarNames = make(map[string]int)
	for n := 0; n < nvars; n++ {
		var name string
		name, err = ncf.inqVarName(n)
		if err != nil {
			return
		}
		ncf.VarNames[name] = n
	}
	return
}

// Type NCfile is the netCDF file object.
type NCfile struct {
	ID              C.int
	Fname           string // file path
	VarNames        map[string]int
	DimNames        map[string]int
	GlobalAttNames  map[string]int
	GlobalAttInt    map[string]int     // int type attributes
	GlobalAttString map[string]string  // string type attributes
	GlobalAttFloat  map[string]float64 // float type attributes
}

// parse errors
func errAnal(err C.int) error {
	var estr error
	if err != 0 {
		strerr := C.nc_strerror(err)
		estr = errors.New(C.GoString(strerr))
	}
	return estr
}

// The function nc_inq_libvers returns a string identifying the version of the 
// netCDF library, and when it was built. 
func InqLibvers() string {
	version := C.nc_inq_libvers()
	return C.GoString(version)
}

//
// The function nc_close closes an open netCDF dataset.
// If the dataset in define mode, nc_enddef will be called before closing. (In this case, if nc_enddef returns an error, nc_abort will automatically be called to restore the dataset to the consistent state before define mode was last entered.) After an open netCDF dataset is closed, its netCDf.ID may be reassigned to the next netCDF dataset that is opened or created.
func (f *NCfile) Close() error {
	e := C.nc_close(f.ID)
	return errAnal(e)
}

// The function nc_sync offers a way to synchronize the disk copy of a netCDF dataset with in-memory buffers
func (f *NCfile) Sync() error {
	e := C.nc_sync(f.ID)
	return errAnal(e)
}

// The function nc_redef puts an open netCDF dataset into define mode, so dimensions, variables, and attributes can be added or renamed and attributes can be deleted. 
func (f *NCfile) ReDef() error {
	e := C.nc_redef(f.ID)
	return errAnal(e)
}

// The function nc_enddef takes an open netCDF dataset out of define mode. The changes made to the netCDF dataset while it was in define mode are checked and committed to disk if no problems occurred. Non-record variables may be initialized to a "fill value" as well. See nc_set_fill. The netCDF dataset is then placed in data mode, so variable data can be read or written.
//It's not necessary to call nc_enddef for netCDF-4 files. With netCDF-4 files, nc_enddef is called when needed by the netcdf-4 library. User calls to nc_enddef for netCDF-4 files still flush the metadata to disk.
//This call may involve copying data under some circumstances. For a more extensive discussion see File Structure and Performance.
//For netCDF-4/HDF5 format files there are some variable settings (the compression, endianness, fletcher32 error correction, and fill value) which must be set (if they are going to be set at all) between the nc_def_var and the next nc_enddef. Once the nc_enddef is called, these settings can no longer be changed for a variable. 
func (f *NCfile) EndDef() error {
	e := C.nc_enddef(f.ID)
	return errAnal(e)
}

// The function nc_def_dim adds a new dimension to an open netCDF dataset in define mode. It returns (as an argument) a dimension ID, given the netCDf.ID, the dimension name, and the dimension length. At most one unlimited length dimension, called the record dimension, may be defined for each classic or 64-bit offset netCDF dataset. NetCDF-4 datasets may have multiple unlimited dimensions. 
func (f NCfile) DefDim(name string, length int) (err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	Clength := C.size_t(length)
	var Cdimidp C.int
	e := C.nc_def_dim(f.ID, cname, Clength, &Cdimidp)
	f.DimNames[name] = int(Cdimidp)
	err = errAnal(e)
	return
}

// The function nc_def_var adds a new variable to an open netCDF dataset in define mode. It returns (as an argument) a variable ID, given the netCDf.ID, the variable name, the variable type, the number of dimensions, and a list of the dimension IDs. 
func (f *NCfile) DefVar(name string, vartype string, dimnames []string) (err error) {
	var cvartype C.nc_type
	switch vartype {
	case "byte":
		cvartype = C.NC_BYTE
	case "char":
		cvartype = C.NC_CHAR
	case "short":
		cvartype = C.NC_SHORT
	case "int":
		cvartype = C.NC_INT
	case "float":
		cvartype = C.NC_FLOAT
	case "double":
		cvartype = C.NC_DOUBLE
	}
	Cvaridp := C.int(0)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ndims := len(dimnames)
	dimids := make([]int, ndims)
	for i := 0; i < ndims; i++ {
		dimids[i] = f.DimNames[dimnames[i]]
	}

	e := C.nc_def_var(f.ID, cname, cvartype, C.int(ndims),
		(*C.int)(unsafe.Pointer(&dimids[0])), &Cvaridp)
	err = errAnal(e)
	f.VarNames[name] = int(Cvaridp)
	return
}

// The function nc_put_vara_ type writes values into a netCDF variable of an open netCDF dataset. The part of the netCDF variable to write is specified by giving a corner and a vector of edge lengths that refer to an array section of the netCDF variable. The values to be written are associated with the netCDF variable by assuming that the last dimension of the netCDF variable varies fastest in the C interface. The netCDF dataset must be in data mode.
// The functions for types ubyte, ushort, uint, longlong, ulonglong, and string are only available for netCDF-4/HDF5 files.
//The nc_put_var() function will write a variable of any type, including user defined type. For this function, the type of the data in memory must match the type of the variable - no data conversion is done.
func (f *NCfile) PutVaraFloat(varname string, start []int64,
	count []int64, fp []float32) error {
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err := errors.New(msg)
		return err
	}
	e := C.nc_put_vara_float(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.float)(unsafe.Pointer(&fp[0])))
	return errAnal(e)
}

func (f *NCfile) PutVaraDouble(varname string, start []int64,
	count []int64, fp []float64) error {
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err := errors.New(msg)
		return err
	}
	e := C.nc_put_vara_double(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.double)(unsafe.Pointer(&fp[0])))
	return errAnal(e)
}

func (f *NCfile) PutVaraText(varname string, start []int64,
	count []int64, val string) error {
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err := errors.New(msg)
		return err
	}
	e := C.nc_put_vara_text(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])), cval)
	return errAnal(e)
}

func (f *NCfile) PutVarFloat(varname string, fp []float32) error {
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err := errors.New(msg)
		return err
	}
	e := C.nc_put_var_float(f.ID, C.int(varid),
		(*C.float)(unsafe.Pointer(&fp[0])))
	return errAnal(e)
}

func (f *NCfile) PutVarDouble(varname string, fp []float64) error {
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err := errors.New(msg)
		return err
	}
	e := C.nc_put_var_double(f.ID, C.int(varid),
		(*C.double)(unsafe.Pointer(&fp[0])))
	return errAnal(e)
}

//The members of the nc_get_vara_ type family of functions read an array of values from a netCDF variable of an open netCDF dataset. The array is specified by giving a corner and a vector of edge lengths. The values are read into consecutive locations with the last dimension varying fastest. The netCDF dataset must be in data mode.
// The functions for types ubyte, ushort, uint, longlong, ulonglong, and string are only available for netCDF-4/HDF5 files.
// The nc_get_vara() function will write a variable of any type, including user defined type. For this function, the type of the data in memory must match the type of the variable - no data conversion is done. 
func (f *NCfile) GetVaraFloat(varname string, start []int64,
	count []int64) (out []float32, err error) {
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err = errors.New(msg)
		return
	}
	length := int64(1)
	for _, dim := range count {
		if dim > 0 {
			length = length * dim
		}
	}
	out = make([]float32, length)
	e := C.nc_get_vara_float(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.float)(unsafe.Pointer(&out[0])))
	err = errAnal(e)
	return
}

func (f *NCfile) GetVaraDouble(varname string, start []int64,
	count []int64) (out []float64, err error) {
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err = errors.New(msg)
		return
	}
	length := int64(1)
	for _, dim := range count {
		if dim > 0 {
			length = length * dim
		}
	}
	out = make([]float64, length)
	e := C.nc_get_vara_double(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.double)(unsafe.Pointer(&out[0])))
	err = errAnal(e)
	return
}

func (f *NCfile) GetVaraText(varname string, start []int64,
	count []int64) (val string, err error) {
	val = cstr
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	varid, ok := f.VarNames[varname]
	if !ok {
		err = fmt.Errorf("Variable %v does not exist in file %v",
			varname, f.Fname)
		return
	}
	e := C.nc_get_vara_text(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])), cval)
	err = errAnal(e)
	val = C.GoString(cval)
	return
}

func (f *NCfile) GetVarDouble(varname string) (
	fp []float64, err error) {
	dims, err := f.VarSize(varname)
	if err != nil {
		return
	}
	n := 1
	for _, dim := range dims {
		n = n * dim
	}
	fp = make([]float64, n)
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err = errors.New(msg)
		return
	}
	e2 := C.nc_get_var_double(f.ID, C.int(varid),
		(*C.double)(unsafe.Pointer(&fp[0])))
	err = errAnal(e2)
	return
}

func (f *NCfile) GetVarFloat(varname string) (
	fp []float32, err error) {
	dims, err := f.VarSize(varname)
	if err != nil {
		return
	}
	n := 1
	for _, dim := range dims {
		n = n * dim
	}
	fp = make([]float32, n)
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err = errors.New(msg)
		return
	}
	e2 := C.nc_get_var_float(f.ID, C.int(varid),
		(*C.float)(unsafe.Pointer(&fp[0])))
	err = errAnal(e2)
	return
}

func (f *NCfile) GetVar1Double(varname string, index []int64) (
	value float64, err error) {
	var cdp C.double
	varid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v", varname, f.Fname)
		err = errors.New(msg)
		return
	}
	e := C.nc_get_var1_double(f.ID, C.int(varid),
		(*C.size_t)(unsafe.Pointer(&index[0])), &cdp)
	err = errAnal(e)
	value = float64(cdp)
	return
}

// The function nc_put_att_ type adds or changes a variable attribute or global attribute of an open netCDF dataset. If this attribute is new, or if the space required to store the attribute is greater than before, the netCDF dataset must be in define mode. 
func (f *NCfile) PutAttText(attName string, varname string, value string) error {
	length := len(value)
	var varid C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err := errors.New(msg)
			return err
		}
		varid = C.int(govarid)
	}
	Ctp := C.CString(value)
	defer C.free(unsafe.Pointer(Ctp))
	cname := C.CString(attName)
	defer C.free(unsafe.Pointer(cname))
	e := C.nc_put_att_text(f.ID, varid,
		cname, C.size_t(length), Ctp)
	return errAnal(e)
}

func (f *NCfile) PutAttInt(attName string, varname string, value int) error {
	length := 1
	var varid C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err := errors.New(msg)
			return err
		}
		varid = C.int(govarid)
	}
	cname := C.CString(attName)
	defer C.free(unsafe.Pointer(cname))
	cval := C.int(value)
	e := C.nc_put_att_int(f.ID, varid,
		cname, C.NC_INT, C.size_t(length), &cval)
	return errAnal(e)
}

func (f *NCfile) inq() (ndims int, nvars int, ngatts int,
	unlimdimid int, err error) {
	var cndims C.int
	var cnvars C.int
	var cngatts C.int
	var cunlimdimid C.int
	e := C.nc_inq(f.ID, &cndims, &cnvars, &cngatts, &cunlimdimid)
	ndims = int(cndims)
	nvars = int(cnvars)
	ngatts = int(cngatts)
	unlimdimid = int(cunlimdimid)
	err = errAnal(e)
	return
}

func (f *NCfile) inqNvars() (nvars int, err error) {
	var cnvars C.int
	e := C.nc_inq_nvars(f.ID, &cnvars)
	err = errAnal(e)
	nvars = int(cnvars)
	return
}

func (f *NCfile) inqVarNatts(varid int) (natts int, err error) {
	var nattsp C.int
	e := C.nc_inq_varnatts(f.ID, C.int(varid), &nattsp)
	natts = int(nattsp)
	err = errAnal(e)
	return
}

func (f *NCfile) inqVarName(varid int) (name string, err error) {
	name = cstr
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	e := C.nc_inq_varname(f.ID, C.int(varid), cname)
	err = errAnal(e)
	name = C.GoString(cname)
	return
}

func (f *NCfile) inqAttName(varid int, attid int) (
	name string, err error) {
	name = cstr
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	e := C.nc_inq_attname(f.ID, C.int(varid), C.int(attid), cname)
	err = errAnal(e)
	name = C.GoString(cname)
	return
}

func (f *NCfile) inqDimName(dimID int) (
	name string, err error) {
	name = cstr
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	e := C.nc_inq_dimname(f.ID, C.int(dimID), cname)
	err = errAnal(e)
	name = C.GoString(cname)
	return
}

// The function nc_copy_att copies an attribute from one open netCDF dataset to another. It can also be used to copy an attribute from one variable to another within the same netCDF.
// If used to copy an attribute of user-defined type, then that user-defined type must already be defined in the target file. In the case of user-defined attributes, enddef/redef is called for ncid_in and ncid_out if they are in define mode. (This is the ensure that all user-defined types are committed to the file(s) before the copy is attempted.) 
// For global attributes, use varname = "global"
func (f *NCfile) CopyAtt(attname string, varname string, fout *NCfile) (
	err error) {
	var varid C.int
	var varidout C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
		varidout = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err = errors.New(msg)
			return
		}
		varid = C.int(govarid)
		govaridout, ok := fout.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, fout.Fname)
			err = errors.New(msg)
			return
		}
		varidout = C.int(govaridout)
	}
	cname := C.CString(attname)
	defer C.free(unsafe.Pointer(cname))
	e := C.nc_copy_att(f.ID, varid, cname, fout.ID, varidout)
	err = errAnal(e)
	return
}

func (f *NCfile) InqAttType(attname string, varname string) (
	attType string, err error) {
	var varid C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err = errors.New(msg)
			return
		}
		varid = C.int(govarid)
	}
	var xtypep C.nc_type
	cname := C.CString(attname)
	defer C.free(unsafe.Pointer(cname))
	e := C.nc_inq_atttype(f.ID, varid, cname, &xtypep)
	err = errAnal(e)
	switch xtypep {
	case C.NC_INT, C.NC_SHORT:
		attType = "int"
	case C.NC_FLOAT, C.NC_DOUBLE:
		attType = "float"
	case C.NC_CHAR:
		attType = "string"
	case C.NC_BYTE:
		attType = "byte"
	default:
		attType = "unknown"
	}
	return
}

// Get attribute value. If variable name ("varname") is "global", return
// global attributes.
func (f *NCfile) GetAttInt(attname string, varname string) (
	val int, err error) {
	var varid C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err = errors.New(msg)
			return
		}
		varid = C.int(govarid)
	}
	cname := C.CString(attname)
	defer C.free(unsafe.Pointer(cname))
	cval := C.int(val)
	e := C.nc_get_att_int(f.ID, varid, cname, &cval)
	err = errAnal(e)
	val = int(cval)
	return
}

func (f *NCfile) GetAttFloat(attname string, varname string) (
	val float64, err error) {
	var varid C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err = errors.New(msg)
			return
		}
		varid = C.int(govarid)
	}
	cname := C.CString(attname)
	defer C.free(unsafe.Pointer(cname))
	cval := C.double(val)
	e := C.nc_get_att_double(f.ID, varid, cname, &cval)
	err = errAnal(e)
	val = float64(cval)
	return
}

func (f *NCfile) GetAttString(attname string, varname string) (
	val string, err error) {
	var varid C.int
	if varname == "global" {
		varid = C.NC_GLOBAL
	} else {
		govarid, ok := f.VarNames[varname]
		if !ok {
			msg := fmt.Sprintf("Variable %v does not exist in file %v",
				varname, f.Fname)
			err = errors.New(msg)
			return
		}
		varid = C.int(govarid)
	}
	cname := C.CString(attname)
	defer C.free(unsafe.Pointer(cname))
	val = cstr
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	e := C.nc_get_att_text(f.ID, varid, cname, cval)
	err = errAnal(e)
	val = strings.Trim(C.GoString(cval), "`")
	return
}

// Function VarSize returns the size of each dimension for the
// specified variable
func (f *NCfile) VarSize(varname string) (dims []int, err error) {
	var ndimsp C.int
	govarid, ok := f.VarNames[varname]
	if !ok {
		msg := fmt.Sprintf("Variable %v does not exist in file %v",
			varname, f.Fname)
		err = errors.New(msg)
		return
	}
	varid := C.int(govarid)
	e := C.nc_inq_varndims(f.ID, varid, &ndimsp)
	if err = errAnal(e); err != nil {
		return
	}
	dimids := make([]int, int(ndimsp))
	dims = make([]int, int(ndimsp))
	e = C.nc_inq_vardimid(f.ID, varid,
		(*C.int)(unsafe.Pointer(&dimids[0])))
	if err = errAnal(e); err != nil {
		return
	}
	var tmplen C.size_t
	for i, id := range dimids {
		e = C.nc_inq_dimlen(f.ID, C.int(id), &tmplen)
		dims[i] = int(tmplen)
		if err = errAnal(e); err != nil {
			return
		}
	}
	return

}