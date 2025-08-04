start-itest-app:
	cd itest/itest_elixir_app && iex --sname itestapp@localhost --cookie 'super_secret' -S mix run

itest-example:
	elixir --sname itest@localhost --cookie 'super_secret' -e 'Node.connect(:"itestapp@localhost");  Node.spawn(:"itestapp@localhost", fn -> send(ItestElixirApp.Worker, "hi") end)'

itest-run-epmd:
	epmd -d

itest-automated:
	cd test_integration && go test -v -run TestIntegrationWithElixirApp

itest-all:
	cd test_integration && go test -v

itest-compile:
	cd test_integration && go test -v -run TestElixirAppCompilation
