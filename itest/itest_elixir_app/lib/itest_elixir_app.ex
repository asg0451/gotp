defmodule ItestElixirApp do
  defmodule Worker do
    use GenServer

    def start_link(arg) do
      GenServer.start_link(__MODULE__, arg, name: __MODULE__)
    end

    # because C genserver idk
    def cast_to_self(msg) do
      GenServer.cast(__MODULE__, msg)
    end
    def call_to_self(msg) do
      GenServer.call(__MODULE__, msg)
    end
    def send_to_self(msg) do
      send(__MODULE__, msg)
    end

    # Callbacks

    @impl true
    def init(stack) do
      {:ok, stack}
    end

    @impl true
    def handle_call(:pop, _from, [head | tail]) do
      {:reply, head, tail}
    end

    @impl true
    def handle_cast({:push, element}, state) do
      {:noreply, [element | state]}
    end

    @impl true
    def handle_info(msg, state) do
      IO.puts("Received message: #{inspect(msg)}")
      {:noreply, state}
    end
  end
end
