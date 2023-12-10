"""
gactar_ccm_activate_trace adds a MemorySubModule to trace activations.

It looks like there's not much information we can get out of the system.
There are no pre-and-post hooks, so we can't get any activation values.

To use it, create one in an ACTR instance:

class foo(ACTR):
    retrieval = Buffer()
    goal = Buffer()
    memory = Memory(retrieval)
    trace = ActivateTrace(memory)
    ...
"""

from python_actr import Chunk, Memory, MemorySubModule


class ActivateTrace(MemorySubModule):
    """
    ActivateTrace provides a method to print a chunk when it is activated.
    """

    def __init__(self, memory: Memory):
        MemorySubModule.__init__(self, memory)

    # We don't have any real info at this point, but we can output
    # the chunk which was activated.
    def activation(self, chunk: Chunk):
        print(f"trace: activated chunk ({chunk})")
        return 0
