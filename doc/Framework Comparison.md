# Framework Comparison

This document is to gather together the comparisons I've made between the three frameworks gactar is working with:

- [python_actr](https://github.com/asmaloney/python_actr) (Python) - a.k.a. **_ccm_**
- [pyactr](https://github.com/jakdot/pyactr) (Python)
- [ACT-R](https://github.com/asmaloney/ACT-R) (Lisp) - a.k.a. **_vanilla_**

When referencing the _ACT-R manual_, I am using the [**ACT-R 7.26+ Reference Manual**](https://github.com/asmaloney/ACT-R/blob/main/docs/reference-manual.pdf) by Dan Bothell from the ACT-R repository.

## Ad-hoc Buffers

According to the ACT-R manual, buffers may only be added through modules (emphasis mine):

> When the **_model_** is defined it may have any number of buffers associated with it. This is the only way to add new buffers to the system – **they do not exist independently of a _module_**. The buffers for a module will provide the interface for other modules to interact with the new module. (pg. 519)

This is a bit confusing because it seems to mix `model` and `module`. If indeed they can only be created by `modules`, then both ccm & pyactr are deviating from that specification.

| framework | supported? | note                                                     |
| --------- | ---------- | -------------------------------------------------------- |
| ccm       | 🟠 **(1)** | added on the model directly                              |
| pyactr    | 🟠 **(2)** | added on the model (e.g. `pyactr_model.set_goal('foo')`) |
| vanilla   | 🟢         | added when calling `define-module`                       |

**(1)** In ccm, you just declare it in the model and start using it like this:

```python
class FooModel(ACTR):
  foo=Buffer()
  ...
```

**(2)** In pyactr, the name of the function `set_goal` is confusing. This means "create a goal-like buffer and set it in the model", not "set the goal to 'foo'".

## Buffer/Module States

| state              | ccm        | pyactr | vanilla    | note                           |
| ------------------ | ---------- | ------ | ---------- | ------------------------------ |
| **Buffer**         |            |        |            | _These query the buffer state_ |
| buffer empty       | 🟢         | 🟢     | 🟢         |                                |
| buffer full        | 🟢         | 🟢     | 🟢         |                                |
| buffer failure     | 🔴         | 🔴     | 🟢         |                                |
| buffer requested   | 🔴         | 🔴     | 🟢         |                                |
| buffer unrequested | 🔴         | 🔴     | 🟢         |                                |
| **Module**         |            |        |            | _These query the module state_ |
| state free         | 🔴         | 🟢     | 🟢         |                                |
| state busy         | 🟠 **(1)** | 🟢     | 🟢         |                                |
| state error        | 🟠 **(1)** | 🟢     | 🟢         |                                |
| error t            | 🔴         | 🔴     | 🟢 **(2)** | _alias_                        |
| error nil          | 🔴         | 🔴     | 🟢 **(3)** | _alias_                        |

**(1)** In ccm, it looks like the state checks are only implemented on some modules. `busy` seems to be on several, but `error` is only on Memory and SOSVision.

**(2)** In vanilla, alias for `state error` (pg. 230).

**(3)** In vanilla, alias for for `– state error` (pg. 230).

## Request Parameters

With the basic buffers, it looks like only `retrieval` has request parameters.

| buffer    | parameter          | value      | ccm | pyactr     | vanilla | note |
| --------- | ------------------ | ---------- | --- | ---------- | ------- | ---- |
| retrieval | recently retrieved | t          | 🔴  | 🔴         | 🟢      |      |
|           |                    | nil        | 🔴  | 🟢 **(1)** | 🟢      |      |
|           |                    | reset      | 🔴  | 🔴         | 🟢      |      |
| retrieval | mp-value           | _(number)_ | 🔴  | 🔴         | 🟢      |      |
| retrieval | rt-value           | _(number)_ | 🔴  | 🔴         | 🟢      |      |

**Note:**

- vanilla handles request parameters declared on any module
- pyactr looks like it only handles retrieval's `recently_retrieved True/False`
- ccm does not seem to support request parameters at all

**(1)** In pyactr, this is applied as separate query, not as part of the request.

## Spreading Activation

Spreading activation is set per-buffer. The defaults in each of the frameworks are different though. In general the default for ccm is **1.0**, pyactr is **0.0**, and vanilla is **0.0** with the exception of the `imaginal` buffer which is **1.0**.

**Note:** I'm not listing all the buffers - just the ones gactar currently supports.

| buffer    | ccm | pyactr      | vanilla | note |
| --------- | --- | ----------- | ------- | ---- |
| goal      | 1.0 | 0.0         | 0.0     |      |
| imaginal  | 1.0 | 0.0 **(1)** | 1.0     |      |
| retrieval | 1.0 | 0.0         | 0.0     |      |

**(1)** pyactr does not have a built-in imaginal module or buffer. The `imaginal` buffer is created on the model like any other ad-hoc buffer. For example:

```python
imaginal = my_model.set_goal(name="imaginal", delay=0.2)
```

## Symbolic vs. Subsymbolic

| framework | symbolic | subsymbolic | note                                                                      |
| --------- | -------- | ----------- | ------------------------------------------------------------------------- |
| ccm       | 🔴       | 🟢          | no apparent way to use symbolic-only                                      |
| pyactr    | 🟢       | 🟢          | controlled by `subsymbolic` parameter set on the model (default: `False`) |
| vanilla   | 🟢       | 🟢          | defaults to symbolic, subsymbolic turned on using `:esc`                  |
