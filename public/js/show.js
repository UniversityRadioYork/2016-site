$(document).ready(function() {
    $("#timeslot-list").hide();
    $("#timeslot-latest").show();

    //Must run before dataTables to apply to all season pages.
    $(".timeslot-filter-season").on( "click", function () {
        filterSeason(this);
    });

    var seasonTable = $("#seasons").DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "pageLength": 5,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
    });
    var timeslotTable = $("#timeslots").DataTable({
        "ordering": false,
        "info":     false,
        "lengthChange": false,
        "sDom": '<"top"i>rt<"bottom"lp><"clear">',
        "pageLength": 5,
        "columnDefs": [
        {
            "targets": [ 0 ],
            "visible": false,
            "searchable": true
        },
        {
            "targets": [ 1,2 ],
            "searchable": false
        },

        ]
    });

    var filterSeason = function(object) {
        $(".timeslot-filter-season").removeClass("active");
        $(object).addClass("active");
        $("#timeslot-latest").hide();
        $("#timeslot-season-title").text("Season " + $(object).find("td:nth-of-type(2)").html())
        timeslotTable
            .columns( 0 )
            .search( $(object).data("seasonid") )
            .draw();

        $("#timeslot-list").show();

        $("#timeslot-list tr").on("click", function () {
            window.location.href = $(this).find("a").attr("href");
        });
    };

    $.urlParam = function (name) {
        var results = new RegExp('[\?&]' + name + '=([^&#]*)').exec(window.location.href);
        if (results == null) {
            return null;
        }
        else {
            return decodeURI(results[1]) || 0;
        }
    }

    var getSeasonFilter = function() {
        if ($.urlParam("seasonID") != null) {
            var seasonTable = $("#seasons").dataTable();
            var rowId = seasonTable.fnFindCellRowIndexes($.urlParam("seasonID"),0);
            $("#seasons").DataTable().row(rowId).show().draw(false);
            filterSeason($("tr[data-seasonid='" + $.urlParam("seasonID") + "']"));
            window.location.href = window.location.href.split('#')[0] + "#timeslot-list";
        }
    }
    getSeasonFilter();
});
