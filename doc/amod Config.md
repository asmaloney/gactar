# amod Config

## gactar Config

Top-level configuration in the `gactar` section.

Example:

```
gactar {
    log_level: 'detail'
    trace_activations: true
}
```

| Config            | Type                                       | Description                                                                              |
| ----------------- | ------------------------------------------ | ---------------------------------------------------------------------------------------- |
| log_level         | string (one of 'min', 'info', or 'detail') | how verbose our logging should be                                                        |
| trace_activations | boolean                                    | output detailed info about activations                                                   |
| random_seed       | positive integer                           | sets the seed to use for generating pseudo-random numbers (allows for reproducible runs) |

## Module Config

gactar supports a handful of modules and configuration options which are set in the `modules` section.

A module's buffer is configured using its name like this:

```
modules {
    memory {
        retrieval { spreading_activation: 0.5 }
    }
}
```

For a list of a specific module's configuration options, run `gactar module info [module name]`.

For a list of _all_ modules and their configuration options, run `gactar module info all`.

Example:

```
modules {
    memory {
        latency_factor: 0.63
        max_spread_strength: 1.6
    }

    goal {
        goal{ spreading_activation: 1.0 }
    }

    extra_buffers {
        foo {}
        bar {}
    }
}
```

### Buffers

All buffers have one configuration option.

| Config               | Type    | Description                                                         |
| -------------------- | ------- | ------------------------------------------------------------------- |
| spreading_activation | decimal | see "Spreading Activation" in "ACT-R 7.26 Reference Manual" pg. 290 |

Mapping defaults is a little messy since they are handled differently in each of the frameworks.

- **ccm** (DMSpreading.weight): defaults to 1.0, however this is incorrect, so we explicitly set it to 0.0
- **pyactr**: if not set explicitly, defaults to 0.0
- **vanilla**: (_see below_)

Here are vanilla's parameters and defaults:

| buffer                 | ACT-R param                 | default |
| ---------------------- | --------------------------- | ------- |
| Aural buffer           | :aural-activation           | 0.0     |
| Aural-location buffer  | :aural-location-activation  | 0.0     |
| Goal buffer            | :ga                         | 0.0     |
| Imaginal buffer        | :imaginal-activation        | 1.0     |
| Imaginal-action buffer | :imaginal-action-activation | 0.0     |
| Manual buffer          | :manual-activation          | 0.0     |
| Production buffer      | :production-activation      | 0.0     |
| Retrieval buffer       | :retrieval-activation       | 0.0     |
| Temporal buffer        | :temporal-activation        | 0.0     |
| Visual buffer          | :visual-activation          | 0.0     |
| Visual-location buffer | :visual-location-activation | 0.0     |
| Vocal buffer           | :vocal-activation           | 0.0     |

Right now, we allow setting the spreading activation on any buffer, however we only generate code for the `goal` buffer.

### Declarative Memory

This is the standard ACT-R declarative memory module.

Module Name: **memory**

Buffer Name: **retrieval**

| Config              | Type    | Description                                                                           | Mapping                                                                                                                     |
| ------------------- | ------- | ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| decay               | decimal | sets the decay for the base-level learning calculation                                | ccm (DMBaseLevel submodule 'decay'): 0.5<br>pyactr (decay) : 0.5<br>vanilla (:bll): nil (recommend 0.5 if used)             |
| finst_size          | integer | how many chunks are retained in memory                                                | ccm (finst_size): 4<br>pyactr (DecMemBuffer.finst): 0<br>vanilla (:declarative-num-finsts): 4                               |
| finst_time          | decimal | how long the finst lasts in memory                                                    | ccm (finst_time): 3.0<br>pyactr: (unsupported? Always ∞ I guess?)<br>vanilla (:declarative-finst-span): 3.0                 |
| instantaneous_noise | decimal | turns on noise calculation & sets instantaneous noise                                 | ccm (DMNoise submodule 'noise')<br>pyactr (instantaneous_noise)<br>vanilla (:ans)                                           |
| latency_exponent    | decimal | latency exponent (f)                                                                  | ccm: (unsupported? Based on the code, it seems to be fixed at 1.0.)<br>pyactr (latency_exponent): 1.0<br>vanilla (:le): 1.0 |
| latency_factor      | decimal | latency factor (F)                                                                    | ccm (latency): 0.05<br>pyactr (latency_factor): 0.1<br>vanilla (:lf): 1.0                                                   |
| max_spread_strength | decimal | turns on the spreading activation calculation & sets the maximum associative strength | ccm (DMSpreading submodule)<br>pyactr (strength_of_association)<br>vanilla (:mas)                                           |
| mismatch_penalty    | decimal | turns on partial matching and sets the penalty in the activation equation to this     | ccm (Partial class)<br>pyactr (partial_matching & mismatch_penalty)<br>vanilla (:mp)                                        |
| retrieval_threshold | decimal | retrieval threshold (τ)                                                               | ccm (threshold): 0.0<br>pyactr (retrieval_threshold): 0.0<br>vanilla (:rt): 0.0                                             |

### Goal

This is the standard ACT-R goal module.

Module Name: **goal**

Buffer Name: **goal**

It does not have any additional configuration options.

### Imaginal

This is the standard ACT-R imaginal module.

Module Name: **imaginal**

Buffer Name: **imaginal**

| Config | Type    | Description                                                     | Mapping                                                                                       |
| ------ | ------- | --------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| delay  | decimal | how long it takes a request to the buffer to complete (seconds) | ccm (ImaginalModule.delay): 0.2<br>pyactr (Goal.delay): 0.2<br>vanilla (:imaginal-delay): 0.2 |

### Procedural

This is the standard ACT-R procedural module.

Module Name: **procedural**

Buffer Name: _none_

| Config              | Type    | Description                                       | Mapping                                                                           |
| ------------------- | ------- | ------------------------------------------------- | --------------------------------------------------------------------------------- |
| default_action_time | decimal | time that it takes to fire a production (seconds) | ccm (production_time): 0.05<br>pyactr (rule_firing): 0.05<br>vanilla (:dat): 0.05 |

### Extra Buffers

This is a gactar-specific module used to add new buffers to the model. According to ACT-R, buffers should only be added through modules, however some implementations allow declaring them wherever you want.

Module Name: **extra_buffers**

Buffer Names: _specified in configuration_

| Config           | Description                                                               |
| ---------------- | ------------------------------------------------------------------------- |
| _buffer name_ {} | the name of the new buffer (with "{}" to allow any buffer config options) |
