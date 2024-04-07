#include <Python.h>

double call_func(PyObject *func, double x, double y)
{
	PyObject *args;
	PyObject *kwargs;
	PyObject *result = 0;
	double retval;

	// Make sure we own the GIL
	PyGILState_STATE state = PyGILState_Ensure();

	// Verify that func is a proper callable
	if (!PyCallable_Check(func)) {
		fprintf(stderr, "call_func: expected a callable\n");
		goto fail;
	}

	args = Py_BuildValue("(dd)", x, y);
	kwargs = NULL;

	result = PyObject_Call(func, args, kwargs);
	Py_DECREF(args);
	Py_XDECREF(kwargs);

	if (PyErr_Occurred()) {
		PyErr_Print();
		goto fail;
	}

	// Verify the result is a float object
	if (!PyFloat_Check(result)) {
		fprintf(stderr, "call_func: callable didn't return a float\n");
		goto fail;
	}

	retval = PyFloat_AsDouble(result);
	Py_DECREF(result);

	PyGILState_Release(state);
	return retval;

fail:
	Py_XDECREF(result);
	PyGILState_Release(state);
	abort();
}

PyObject *import_name(const char *modname, const char *symbol)
{
	PyObject *u_name, *module;
	u_name = PyUnicode_FromString(modname);
	module = PyImport_Import(u_name);
	Py_DECREF(u_name);

	return PyObject_GetAttrString(module, symbol);
}

int main2()
{
    PyObject * pow_func;
    double x;
    Py_Initialize();

    // Get a reference to the math.pow function
    pow_func = import_name("math", "pow");

    // Call it using our call_func() code
    for (x = 0.0; x < 10.0; x += 1.0) {
        printf("% 0.2f % 0.2f\n", x, call_func(pow_func, x, 2.0));
    }

    Py_DECREF(pow_func);
    Py_Finalize();
    return 0;
}
