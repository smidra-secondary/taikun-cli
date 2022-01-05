Context 'billing/rule/delete'

  setup() {
    name=$(_rnd_name)
    cname=$(_rnd_name)
    cid=$(taikun billing credential create -p $PROMETHEUS_PASSWORD -u $PROMETHEUS_URL -l $PROMETHEUS_USERNAME $cname -I)

    flags="-b $cid -l foo=bar -m foo --price 1 --price-rate 5 -t count"
    id=$(taikun billing rule create $name $flags -I)
  }

  cleanup() {
    taikun billing rule delete $id -q 2>/dev/null || true
    taikun billing credential delete $cid -q 2>/dev/null || true
  }

  BeforeEach 'setup'
  AfterEach 'cleanup'

  Example 'delete nonexistent billing rule'
    When call taikun billing rule delete 0
    The status should equal 1
    The stderr should include 404
    The stderr should include 'Error: Failed to delete one or more resources'
  End

  Example 'delete existing billing rule'
    When call taikun billing rule delete $id
    The status should equal 0
    The output should include 'was deleted successfully'
    The output should include "$id"
    The lines of output should equal 1
  End

  Example 'delete existing and nonexistent billing rules'
    When call taikun billing rule delete 0 $id
    The status should equal 1
    The output should include 'was deleted successfully'
    The output should include "$id"
    The stderr should include 404
    The stderr should include 'Error: Failed to delete one or more resources'
  End

End
