defmodule ItestElixirAppTest do
  use ExUnit.Case
  doctest ItestElixirApp

  test "greets the world" do
    assert ItestElixirApp.hello() == :world
  end
end
